/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "读取数据库结构生成实体",
	Long:  `读取数据库结构生成实体`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gen called")
		fmt.Println(dbtype)
		fmt.Println(database)
		fmt.Println(host)
		fmt.Println(passwd)
		fmt.Println(table)
	},
}
var (
	database string
	dbtype   string
	host     string
	passwd   string
	table    []string
)

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	genCmd.Flags().StringVar(&dbtype, "dbtype", "mysql", "数据库类型")
	genCmd.Flags().StringVar(&database, "db", "user", "数据名称")
	genCmd.Flags().StringVar(&host, "host", "localhost:3306", "数据库地址")
	genCmd.Flags().StringVar(&host, "passwd", "123456", "数据库密码")
	genCmd.Flags().StringSliceVar(&table, "table", nil, "指定要生成的表")
}
