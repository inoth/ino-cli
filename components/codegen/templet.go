package codegen

import (
	"os"
	"text/template"

	"fmt"
	"unicode"

	"strings"
	"sync"

	"github.com/inoth/ino-cli/config"
	"github.com/inoth/ino-cli/db/model"
	"github.com/inoth/ino-cli/util/inofile"
)

var (
	fileTemplate map[string]IHandler
	dataTypeMap  map[string]string
	projectName  string
	dbName       string
	useDbName    bool = false

	Model = "package model\n" +
		"type {{.TableName}} struct {\n" +
		"{{range .Fields}}\n" +
		"{{.Field}} {{.DbType}} `gorm:\"{{.DbField}}\"`" +
		"{{end}}\n}"
)

func InitTemplet() {
	dataTypeMap = make(map[string]string)
	dataTypeMap["varchar"] = "string"
	dataTypeMap["decimal"] = "decimal"
	dataTypeMap["int"] = "int"
	dataTypeMap["tinyint"] = "int"
	dataTypeMap["datetime"] = "time.Time"

	fileTemplate = make(map[string]IHandler)
	fileTemplate["Model"] = &EntityHandle{
		Folder: "Model",
		Tmpl:   Model}

	dbName = config.Cfg.GetString("dbName")
	projectName = config.Cfg.GetString("project")
}

type TmplData struct {
	ProjectName string
	TableName   string
	DbTableName string
	UseDbName   bool
	Fields      []Field
}
type Field struct {
	DbField string
	Field   string
	Key     string
	IsNull  string
	Desc    string
	DbType  string
}

type IHandler interface {
	Process(wg *sync.WaitGroup, tableName string)
}

type NormalHandle struct {
	Folder string
	Tmpl   string
}

func (e *NormalHandle) Process(wg *sync.WaitGroup, tableName string) {
	defer wg.Done()
	tmpl, err := template.New(tableName).Parse(e.Tmpl)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	data := &TmplData{
		ProjectName: projectName,
		TableName:   NameHandler(tableName),
		DbTableName: tableName,
		UseDbName:   useDbName,
	}
	path := fmt.Sprintf("./%v/%v/%v.go", projectName, e.Folder, tableName)
	// fmt.Printf("写入： %v", path)
	err = inofile.CreateFileBytes(path, func(f *os.File) error {
		e := tmpl.Execute(f, data)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

type EntityHandle struct {
	Folder string
	Tmpl   string
}

func (e *EntityHandle) Process(wg *sync.WaitGroup, tableName string) {
	defer wg.Done()

	cols := model.GetColumns(dbName, tableName)
	if len(cols) <= 0 {
		fmt.Printf("%v: 未找到有效列\n", tableName)

		return
	}
	tmpl, err := template.New(tableName).Parse(e.Tmpl)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	data := &TmplData{
		ProjectName: projectName,
		TableName:   NameHandler(tableName),
		DbTableName: tableName,
		UseDbName:   useDbName,
	}
	data.Fields = make([]Field, len(cols))
	for i, col := range cols {
		data.Fields[i] = Field{DbType: MatchType(col.DataType),
			DbField: col.ColName,
			Field:   NameHandler(col.ColName),
			Desc:    col.ColDesc}
	}
	path := fmt.Sprintf("./%v/%v/%v.go", projectName, e.Folder, tableName)
	// fmt.Printf("写入： %v", path)
	err = inofile.CreateFileBytes(path, func(f *os.File) error {
		e := tmpl.Execute(f, data)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func MatchType(dbType string) string {
	if val, ok := dataTypeMap[dbType]; ok {
		return val
	} else {
		return "string"
	}
}

func NameHandler(name string) string {
	var tmp string
	strs := strings.Split(name, "_")
	for _, str := range strs {
		if len(str) <= 0 {
			continue
		}
		b := []byte(str)
		b[0] = byte(unicode.ToUpper(rune(b[0])))
		tmp += string(b)
	}
	return tmp
}

func SetTmplGlobleVal(prjName, dbnm string, fmtDbName int) {
	projectName = prjName
	dbName = dbnm
	useDbName = (fmtDbName <= 0)
}
