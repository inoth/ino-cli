package initialize

import (
	"archive/zip"
	"bytes"
	"ino-cli/inocmd"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	emptyProjectName = "ino-empty"
	zipPath          = "https://github.com/inoth/ino-empty/archive/refs/heads/main.zip"
)

var (
	DefaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
)

func Help() {
	log.Fatalln(strings.TrimLeft(`
USAGE    
    ino init NAME
ARGUMENT 
    NAME  name for the project. It will create a folder with NAME in current directory.
          The NAME will also be the module name for the project.
EXAMPLES
    ino init my-app
    ino init my-project-name`, DefaultTrimChars))
}

func Run() {
	projectName := inocmd.GetArg(2)
	if projectName == "" {
		projectName = emptyProjectName
	}
	respData, err := http.Get(zipPath)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer respData.Body.Close()
	zipData, err := ioutil.ReadAll(respData.Body)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if !unZip(zipData, projectName) {
		log.Fatal("fail")
		return
	}

	log.Fatal("finish")
}

func unZip(data []byte, dst string) bool {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	pt := ""
	if dst != "" {
		pt = emptyProjectName + "-main"
		if err = os.MkdirAll(dst, 0755); err != nil {
			log.Fatal(err.Error())
			return false
		}
	}
	for _, file := range zr.File {
		path := filepath.Join(dst, file.Name)
		if pt != "" {
			path = strings.Replace(path, pt, "", -1)
		}
		println(path)
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(path, file.Mode()); err != nil {
				log.Fatal(err.Error())
				return false
			}
			continue
		}

		dir := filepath.Dir(path)
		if len(dir) > 0 {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return false
				}
			}
		}

		fr, err := file.Open()
		if err != nil {
			log.Fatal(err.Error())
			return false
		}
		defer fr.Close()

		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			log.Fatal(err.Error())
			return false
		}
		defer fw.Close()

		_, err = io.Copy(fw, fr)
		if err != nil {
			log.Fatal(err.Error())
			return false
		}
	}
	return true
}
