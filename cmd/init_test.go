package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/inoth/ino-cli/cmd"
	"github.com/inoth/ino-cli/codegen"
	"github.com/inoth/ino-cli/config"
	"github.com/inoth/ino-cli/db"
	"github.com/inoth/ino-cli/global"
	"github.com/inoth/ino-cli/util"
)

func TestDownload(t *testing.T) {
	var (
		zipPath = "https://github.com/inoth/ino-empty/archive/refs/heads/v2.zip"
	)

	fileName, err := cmd.DownloadPackage(zipPath)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("ok; %v", fileName)
}

func TestUnZipAndOutPut(t *testing.T) {
	var (
		// path = "/Users/inoth/self data/work data/github/inoth/ino-cli/remote/zip/ino-empty-2.zip"
		path = "remote/zip/ino-empty-2.zip"
	)

	err := cmd.UnZipAndOutput(path)
	if err != nil {
		t.Errorf(err.Error())
	}
	// cur, _ := os.Executable()
	// newPath := filepath.Dir(cur)
	// t.Logf("%v", cur)
	// t.Logf("%v", newPath)
	t.Logf("ok;")
}

func TestPath(t *testing.T) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// 可执行文件的路径
	fmt.Println(ex)

	//	获取执行文件所在目录
	exPath := filepath.Dir(ex)
	fmt.Println("可执行文件路径 :" + exPath)

	// 使用EvalSymlinks获取真是路径
	realPath, err := filepath.EvalSymlinks(exPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("符号链接真实路径:" + realPath)
}

func TestGenEntity(t *testing.T) {
	dbstr := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", "glasses_mp", "yKxm5jhhwqslOrIj", "rm-uf64ywy4i93cxvt9tzo.mysql.rds.aliyuncs.com:3306", "zk_board")
	util.Must(global.Register(
		config.ViperConfig{}.SetDefaultValue(map[string]interface{}{
			"db":      dbstr,
			"dbName":  "glasses_mp",
			"tables":  nil,
			"project": "testproject",
		}),
		&db.MysqlConnect{}).
		Init().
		Run(&codegen.MysqlGormEntity{}))
}
