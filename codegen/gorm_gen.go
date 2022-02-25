package codegen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/inoth/ino-cli/config"
	"github.com/inoth/ino-cli/db/model"
	"github.com/inoth/ino-cli/util/inofile"
)

var (
	handleLimit = 100
)

type MysqlGormEntity struct{}

func (MysqlGormEntity) ServeStart() error {
	var tables []model.TableInfo
	tableName := config.Cfg.GetStringSlice("tables")
	if len(tableName) <= 0 {
		tables = model.GetTables(config.Cfg.GetString("dbName"))
	} else {
		for _, tbl := range tableName {
			if len(tbl) <= 0 {
				continue
			}
			tables = append(tables, model.TableInfo{TableName: tbl})
		}
	}

	tn := len(tables)
	fmt.Printf("开始处理，一共%d个表\n", tn)
	if tn <= 0 {
		return errors.New("没有找到需要生成表")
	}
	InitTemplet()
	// err := os.MkdirAll(config.Cfg.GetString("project"), 0755)
	// if err != nil {
	// 	return err
	// }
	CreateFolder(config.Cfg.GetString("project"))
	ch_progress := make(chan string, tn)
	// 并行处理最大限制
	ch_limit := make(chan struct{}, handleLimit)
	for _, table := range tables {
		ch_limit <- struct{}{}
		fmt.Printf("开始处理%v\n", table.TableName)
		go func(ctx context.Context, progress chan string, limit chan struct{}, table model.TableInfo) {
			HandlerTable(table)
			ch_progress <- fmt.Sprintf("%v 处理完成", table.TableName)
			<-ch_limit
			return
		}(context.TODO(), ch_progress, ch_limit, table)
	}

	curProgress := 0
	for curProgress < tn {
		r := <-ch_progress
		fmt.Println(r)
		curProgress++
	}
	return nil
}

func CreateFolder(projectName string) {
	for k, _ := range fileTemplate {
		path := fmt.Sprintf("./%v/%v", projectName, k)
		if inofile.PathExists(path) {
			continue
		}
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	}
}

func HandlerTable(table model.TableInfo) {
	wg := sync.WaitGroup{}
	for _, temp := range fileTemplate {
		wg.Add(1)
		go temp.Process(&wg, table.TableName)
	}
	wg.Wait()
}
