package main

import "fmt"

const (
	tableNameSql = "SELECT TABLE_TYPE, TABLE_NAME,TABLE_COMMENT FROM TABLES WHERE table_schema = 'mqjx_base'"
	tableInfoSql = `SELECT COLUMN_NAME,DATA_TYPE, COLUMN_COMMENT,COLUMN_DEFAULT,COLUMN_KEY,COLUMN_TYPE,EXTRA FROM COLUMNS WHERE TABLE_NAME = '%s'  and TABLE_SCHEMA ='mqjx_base'`
)

func main() {
	TouchDir()
	InitDB()
	tables, err := GetDataMap(BaseInfo, tableNameSql)
	if err != nil {
		fmt.Println("table info errorL:", err.Error())
	}

	var list []string
	for _, v := range tables {
		list = append(list, v["TABLE_NAME"].(string))
	}

	RanderLoad(list)
	RanderTemplate(list)
	fmt.Println("生成结束")
}
