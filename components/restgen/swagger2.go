package restgen

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/inoth/ino-cli/util/inofile"
)

func init() {
	add("2.0", func() SwagParsingProcess {
		return &Swag2Process{HandlerLimit: 100}
	})
}

type Swag2Process struct {
	HandlerLimit int
	doc          *SwaggerDocs
}

func (s *Swag2Process) SetContent(doc *SwaggerDocs) error {
	s.doc = doc
	return nil
}

// 获取文档基础信息
// 开启协程给每个path并行生成内容，塞入同一个channel中
// 往同一文件中写入内容，最后生成到磁盘上
func (s *Swag2Process) Process() error {
	ch_pathContent := make(chan string)
	defer close(ch_pathContent)
	ch_limit := make(chan struct{}, s.HandlerLimit)
	defer close(ch_pathContent)
	for url, method := range s.doc.Paths {
		ch_limit <- struct{}{}
		go func(host, url string, method map[string]RequestBlock, ch_pathContent chan string) {
			// 渲染模版
			fmt.Printf("开始生成[%s]\n", url)
			// 统合生成结果
			r := s.process(host, url, method)
			// fmt.Printf("[%s]生成结果：%s\n", url, r)
			<-ch_limit // 释放占用
			ch_pathContent <- r
			return
		}(s.doc.Host, url, method, ch_pathContent)
	}
	fmt.Println("开始拼接处理返回...")
	var result bytes.Buffer
	for i := 0; i < len(s.doc.Paths); i++ {
		result.WriteString(<-ch_pathContent)
	}
	return inofile.WriteToFile(result.Bytes(), fmt.Sprintf("%s.http", s.doc.Info.Title), os.ModePerm)
}

func (s *Swag2Process) process(host, url string, method map[string]RequestBlock) string {
	var res bytes.Buffer
	for met, reqBody := range method {
		temp := &RestModule{
			Method:   met,
			Summary:  reqBody.Summary,
			Host:     host,
			Url:      url,
			ShowBody: true,
		}
		if reqBody.Consumes != nil && len(reqBody.Consumes) > 0 {
			temp.Consumes = reqBody.Consumes[0]
		}
		if met == "get" {
			temp.ShowBody = false
		}
		for i := 0; i < len(reqBody.Parameters); i++ {
			if reqBody.Parameters[i].In == "path" {
				// path 参数无需处理
			} else if reqBody.Parameters[i].In == "query" {
				// 处理get请求参数
				if i == 0 {
					temp.Url += fmt.Sprintf("?")
				}
				temp.Url += fmt.Sprintf("%s=%s&", reqBody.Parameters[i].Name, reqBody.Parameters[i].Desc)
			} else {
				// 添加到body中
				if temp.Params == nil {
					temp.Params = make(map[string]interface{})
				}
				if reqBody.Parameters[i].Schema != (ParamBlockSchema{}) {
					// 判断直接引用，或者是数组中引用
					if reqBody.Parameters[i].Schema.Ref == "" {
						// 数组中引用处理
						temp.Params[reqBody.Parameters[i].Name] = s.refBodyHandler(reqBody.Parameters[i].Schema.Items.Ref)
					}
					// 直接查找引用
					temp.Params[reqBody.Parameters[i].Name] = s.refBodyHandler(reqBody.Parameters[i].Schema.Ref)
				} else { // 直接添加
					temp.Params[reqBody.Parameters[i].Name] = modifyingParam(reqBody.Parameters[i])
				}
			}
		}
		// // 根据模版套入内容生成，返回
		// tmpl, err := template.New(temp.Summary).Parse(RestHttpTmpl)
		// if err != nil {
		// 	fmt.Printf("%v\n", err)
		// 	return ""
		// }
		// buf := bytes.Buffer{}
		// err = tmpl.Execute(&buf, temp)
		// if err != nil {
		// 	fmt.Printf("%v\n", err)
		// 	return ""
		// }
		res.Write(genContent(temp.Summary, RestHttpTmpl, temp))
	}
	return res.String()
}

func (s *Swag2Process) refBodyHandler(refBody string) interface{} {
	tmp := strings.Split(refBody, "/")
	bodyName := tmp[len(tmp)-1]
	body, ok := s.doc.Definitions[bodyName]
	if !ok { // 拿到引用的model内容
		return "\"\""
	}
	bodyTmpl := &BodyModule{
		IsArray: false,
		Body:    map[string]interface{}{},
	}
	for propName, prop := range body.Prop {
		// 遍历，判断是不是引用
		if prop.Ref != "" {
			// 直接跳引用
			bodyTmpl.Body[propName] = s.refBodyHandler(prop.Ref)
		} else if prop.Type == "array" {
			if prop.Items.Ref == "" {
				bodyTmpl.Body[propName] = "[]"
			} else {
				// 条数组中的引用
				bodyTmpl.IsArray = true
				bodyTmpl.Body[propName] = s.refBodyHandler(prop.Items.Ref)
			}
		} else {
			bodyTmpl.Body[propName] = fmt.Sprintf("\"%s\"", prop.Type)
		}
	}
	return string(genContent(bodyName, BodyTmpl, bodyTmpl))
}

// 对生成参数粗剪
func modifyingParam(param ParamBlock) interface{} {
	switch param.Type {
	case "string":
		// if param.Example != "" {
		// 	return fmt.Sprintf("\"%s\"", param.Example)
		// }
		return fmt.Sprintf("\"%s\"", param.Type)
	case "integer":
		// if param.Example != "" {
		// 	return param.Example
		// }
		return fmt.Sprintf("\"%s\"", param.Type)
	case "int32":
		// if param.Example != "" {
		// 	return param.Example
		// }
		return fmt.Sprintf("\"%s\"", param.Type)
	case "int64":
		// if param.Example != "" {
		// 	return param.Example
		// }
		return fmt.Sprintf("\"%s\"", param.Type)
	default:
		return fmt.Sprintf("\"%s\"", param.Type)

	}
}

func genContent(name, tmplstr string, data interface{}) []byte {
	// 根据模版套入内容生成，返回
	tmpl, err := template.New(name).Parse(tmplstr)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}
	return buf.Bytes()
}
