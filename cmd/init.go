/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/inoth/ino-cli/util"
	"github.com/spf13/cobra"
)

var (
	templateName  = "ino-empty-2.zip"
	projectName   = "defaultProject"
	zipPath       = "https://github.com/inoth/ino-empty/archive/refs/heads/v2.zip"
	remotePackage = true
	downloadPath  = "remote/zip/"
	exportPath    = "export"
)
var (
	p_zipPath string
	p_remote  string
	p_output  string
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
		if p_remote == "no" {
			remotePackage = false
		}
		if p_zipPath != "" {
			zipPath = p_zipPath
		}
		if p_output != "" {
			exportPath = fmt.Sprintf("./%s/", p_output)
		}
		fmt.Printf("%v initialization...\n", projectName)
		util.Must(InitProject())
	},
}

func init() {
	initCmd.Flags().StringVar(&p_zipPath, "zip", "", "模版地址")
	initCmd.Flags().StringVar(&p_remote, "remote", "yes", "远程包：yes/no")
	initCmd.Flags().StringVar(&p_output, "output", "export", "生成输出地址")
	rootCmd.AddCommand(initCmd)
}

// 获取远程包前提下，监测本地下载是否存在，否则去拉取远程包
// 解压缩
// 替换项目名称
// 写入文件夹
func InitProject() error {
	var fileName string
	if remotePackage {
		var ok bool
		fileName, ok = checkLocalFile()
		if !ok {
			var err error
			fileName, err = downloadPackage()
			if err != nil {
				return err
			}
		}
	}
	// println(fileName)
	// return nil
	err := unZipAndOutput(downloadPath + fileName)
	// fmt.Printf("%v", err)
	return err
}

// 下载初始化资源包到本地
func downloadPackage() (string, error) {
	fmt.Println("download package...")
	respData, err := http.Get(zipPath)
	if err != nil {
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
	fileName := templateName
	fmt.Println("package write to disk..." + downloadPath + fileName)
	util.Makedir(downloadPath, 0755)
	err = util.WriteToFile(zipData, downloadPath+fileName, 0755)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func checkLocalFile() (string, bool) {
	fileNames, err := util.GetFileNames(downloadPath)
	if err != nil {
		return "", false
	}
	file := templateName
	for _, fileName := range fileNames {
		if fileName == file {
			return fileName, true
		}
	}
	return "", false
}

// 解压后资源放入内存中
func unZipAndOutput(path string) error {
	fmt.Printf("unzip %s package...\n", path)
	util.Makedir(exportPath, 0755)

	// buf, err := util.ReadFile(path)
	// if err != nil {
	// 	return err
	// }

	// zr, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	zr, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	err = util.Makedir(projectName, 0755)
	if err != nil {
		return err
	}
	for _, item := range zr.File {
		path := filepath.Join(projectName, item.Name)
		if item.FileInfo().IsDir() {
			err = util.Makedir(projectName, 0755)
			if err != nil {
				return err
			}
			continue
		}
		dir := filepath.Dir(path)
		if len(dir) > 0 {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
			}
		}
		// 写入磁盘
		fmt.Println("output...")
		fr, err := item.Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		content, _ := ioutil.ReadAll(fr)
		newContent := strings.Replace(string(content), "defaultProject", projectName, -1)
		err = util.WriteToFileByString(newContent, path, item.Mode())
		if err != nil {
			return err
		}

	}
	return nil
}

// // 输出成品到文件目录
// func export(data []byte, path string) error {
// 	return nil
// }
