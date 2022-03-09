package restgen

import (
	"encoding/json"
	"errors"
	"fmt"
)

type RestModule struct {
	Summary  string
	Method   string
	Host     string
	Url      string
	Consumes string
	Params   map[string]string
}

var RestHttpTmpl = `
### {{.Summary}}

{{.Method}} http://{{.Host}}{{.Url}} HTTP/1.1
content-type: {{.Consumes}}

{
   {{- range $key, $value := .Params }}
   "{{$key}}": "{{$value}}",
   {{- end}}
}

`

type SwagParsingProcess interface {
	SetContent(doc *SwaggerDocs) error // 导入解析数据
	Process() error                    // 处理内容
}

type Creator func() SwagParsingProcess

var Process = map[string]Creator{}

func add(name string, creator Creator) {
	Process[name] = creator
}

func GameStart(fileType string, data []byte) error {
	var doc SwaggerDocs
	fileType = "json" // 暂时默认直接json格式
	switch fileType {
	case "json":
		doc = doParseJson(data)
	case "yaml":
		doc = doParseYaml(data)
	default:
		return errors.New("Invalid file type.")
	}
	var process SwagParsingProcess
	switch doc.Swagger {
	case "2.0":
		if tmp, ok := Process["2.0"]; ok {
			process = tmp()
		} else {
			return errors.New("Invalid swagger version.")
		}
	case "3.0":
		if tmp, ok := Process["3.0"]; ok {
			process = tmp()
		} else {
			return errors.New("Invalid swagger version.")
		}
	default:
		return errors.New("Invalid swagger version.")
	}
	err := process.SetContent(&doc)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	err = process.Process()
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}

func doParseJson(data []byte) SwaggerDocs {
	var doc SwaggerDocs
	err := json.Unmarshal(data, &doc)
	if err != nil {
		panic(errors.New("Unable to parse content."))
	}
	return doc
}

func doParseYaml(data []byte) SwaggerDocs {
	panic(errors.New("Not supported at this time."))
}

type SwaggerDocs struct {
	Swagger string `json:"swagger"`
	Info    struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"info"`
	Host  string                             `json:"host"`
	Paths map[string]map[string]RequestBlock `json:"paths"`
}

type RequestBlock struct {
	Summary    string       `json:"summary"`
	Consumes   []string     `json:"consumes"`
	Parameters []ParamBlock `json:"parameters"`
}

type ParamBlock struct {
	Name string `json:"name"`
	In   string `json:"in"`
	Desc string `json:"description"`
}
