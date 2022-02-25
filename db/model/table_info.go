package model

import "github.com/inoth/ino-cli/db"

type TableInfo struct {
	TableName string
	TableDesc string
}
type ColumnsInfo struct {
	ColName  string
	DataType string
	ColDesc  string
}

func GetTables(dbName string) []TableInfo {
	var tables []TableInfo
	db.DB.Raw("SELECT TABLE_NAME as TableName,TABLE_COMMENT as TableDesc FROM INFORMATION_SCHEMA.`TABLES` WHERE TABLE_SCHEMA = ?", dbName).Scan(&tables)
	return tables
}

func GetColumns(dbName, tableName string) []ColumnsInfo {
	var cols []ColumnsInfo
	db.DB.Raw("SELECT COLUMN_NAME as ColName,DATA_TYPE as DataType,COLUMN_COMMENT as ColDesc,IS_NULLABLE AS `IsNull`,COLUMN_KEY AS `Key` FROM INFORMATION_SCHEMA.`COLUMNS` WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", dbName, tableName).Scan(&cols)
	return cols
}
