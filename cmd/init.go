/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/inoth/ino-cli/util"
	"github.com/spf13/cobra"
)

var (
	projectName   = "defaultProject"
	zipPath       = "https://github.com/inoth/ino-empty/archive/refs/heads/v2.zip"
	remotePackage = true
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [Project]",
	Short: "初始化项目手脚架",
	Long:  `初始化项目手脚架`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if name != "" {
			projectName = name
		}
		fmt.Printf("%v initialization...\n", projectName)
		util.Must(InitProject())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// 获取远程包前提下，监测本地下载是否存在，否则去拉取远程包
// 解压缩
// 替换项目名称
// 写入文件夹
func InitProject() error {

	return nil
}

// 下载初始化资源包到本地
func downloadPackage() (string, error) {
	fmt.Println("download package...")
	respData, err := http.Get(zipPath)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}
	defer respData.Body.Close()
	zipData, err := ioutil.ReadAll(respData.Body)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}
	// 写到磁盘上，回传路径
	// return zipData, nil
	return string(zipData), nil
}

// 解压后资源放入内存中
func unZip(path string) ([]byte, error) {
	return nil, nil
}

// 输出成品到文件目录
func export(data []byte, path string) error {
	return nil
}
