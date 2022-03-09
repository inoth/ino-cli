/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/inoth/ino-cli/components/codegen"
	"github.com/inoth/ino-cli/config"
	"github.com/inoth/ino-cli/db"
	"github.com/inoth/ino-cli/global"
	"github.com/inoth/ino-cli/util"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen [Project]",
	Short: "读取数据库结构生成实体",
	Long:  `读取数据库结构生成实体`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		if projectName == "" {
			fmt.Println("未设置项目包名称")
			return
		}
		runCodegen(projectName)
	},
}
var (
	database string
	// dbtype   string
	host    string
	account string
	passwd  string
	table   []string
)

func init() {
	// genCmd.Flags().StringVar(&dbtype, "dbtype", "mysql", "数据库类型")
	// genCmd.Flags().StringVar(&projectName, "table", "defaultProject", "项目名称")
	genCmd.Flags().StringVar(&database, "db", "dbname", "数据库名称")
	genCmd.Flags().StringVar(&host, "host", "localhost:3306", "数据库地址和端口")
	genCmd.Flags().StringVar(&account, "user", "user", "数据库账号")
	genCmd.Flags().StringVar(&passwd, "passwd", "123456", "数据库密码")
	genCmd.Flags().StringSliceVar(&table, "table", nil, "指定要生成的表")
	rootCmd.AddCommand(genCmd)
}

func runCodegen(name string) {
	dbstr := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", account, passwd, host, database)
	util.Must(global.Register(
		config.ViperConfig{}.SetDefaultValue(map[string]interface{}{
			"db":      dbstr,
			"dbName":  database,
			"tables":  table,
			"project": name,
		}),
		&db.MysqlConnect{}).
		Init().
		Run(&codegen.MysqlGormEntity{}))
}
