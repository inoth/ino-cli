package inofile

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetBasePath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	exPath := filepath.Dir(ex)
	realPath, err := filepath.EvalSymlinks(exPath)
	if err != nil {
		return ""
	}
	return realPath
}

func GetFileNames(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	fileNames := make([]string, len(files))
	for i := 0; i < len(files); i++ {
		fileNames[i] = files[i].Name()
	}
	return fileNames, nil
}

func ReadFile(path string) ([]byte, error) {
	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer fw.Close()
	buf, err := ioutil.ReadAll(fw)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func WriteToFile(data []byte, path string, fileMode fs.FileMode) error {
	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer fw.Close()
	n, err := fw.Write(data)
	if err != nil {
		return err
	}
	fmt.Printf("下载zip，大小%d\n", n)
	return nil
}

func WriteToFileByString(data string, path string, fileMode fs.FileMode) error {
	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer fw.Close()
	_, err = fw.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}
