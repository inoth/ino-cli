/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/inoth/ino-cli/components/restgen"
	"github.com/inoth/ino-cli/util"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "解析swagger文档，生成http文件",
	Long:  `解析读取swagger文档内容，解析生成http请求测试文件`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rest client file gen start...")
		util.Must(RunHttpGen())
	},
}

var (
	target      string // 目标类型
	swaggerPath string // 文档所在地址
)

func init() {
	rootCmd.AddCommand(restCmd)
	restCmd.Flags().StringVar(&target, "target", "http", "文档目标地址，默认http获取;(file获取暂定)")
	restCmd.Flags().StringVar(&swaggerPath, "path", "http://localhost:8080/swagger/api-docs.json", "文件所在地址")
}

// 获取文件内容
// 判别文件内容解析内容
// 获取文档基础信息
// 开启协程给每个path并行生成内容，塞入同一个channel中
// 往同一文件中写入内容，最后生成到磁盘上
func RunHttpGen() (err error) {
	var fileContent []byte
	switch target {
	case "file":
		fileContent, err = doFileGet()
	case "http":
		fileContent, err = doHttpGet()
	default:
		err = errors.New("Invalid file source.")
	}
	if err != nil {
		return
	}
	// 获取filetype一并交由 httpgen 模块处理
	err = restgen.GameStart(checkTargetType(), fileContent)
	return
}

func checkTargetType() string {
	tmp := strings.Split(swaggerPath, ".")
	return tmp[len(tmp)-1]
}

func doFileGet() ([]byte, error) {
	return nil, errors.New("Not supported at this time.")
}

func doHttpGet() ([]byte, error) {
	fmt.Println("download document...")
	respData, err := http.Get(swaggerPath)
	if err != nil {
		return nil, err
	}
	defer respData.Body.Close()
	data, err := ioutil.ReadAll(respData.Body)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return data, nil
}
