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

	"github.com/inoth/ino-cli/util"
	"github.com/spf13/cobra"
)

const defaultProject = "defaultProject"
const zipPath = ""

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [Project]",
	Short: "初始化项目手脚架",
	Long:  `初始化项目手脚架`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		if projectName == "" {
			util.Must(errors.New("projectName invalid."))
		}
		fmt.Printf("%v initialization...\n", projectName)
		util.Must(InitProject(projectName))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// 拉取远程资源包
// 解压缩
// 替换项目名称
// 写入文件夹
func InitProject(projectName string) error {
	return nil
}

// 下载初始化资源包
func DownloadPackage() ([]byte, error) {
	fmt.Println("download package...")
	respData, err := http.Get(zipPath)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer respData.Body.Close()
	zipData, err := ioutil.ReadAll(respData.Body)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return zipData, nil
}
