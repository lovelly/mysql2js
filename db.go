package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const driverName = "mysql"

var (
	once     = &sync.Once{}
	BaseDB   *sqlx.DB
	BaseInfo *sqlx.DB
)

var Server struct {
	BaseDbHost     string
	BaseDbPort     int
	BaseDbName     string
	BaseDbUsername string
	BaseDbPassword string
}

func GetBaseDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.BaseDbUsername, Server.BaseDbPassword, Server.BaseDbHost, Server.BaseDbPort, Server.BaseDbName, "parseTime=true&loc=Local&interpolateParams=true")
	return s
}
func GetBaseInfoDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.BaseDbUsername, Server.BaseDbPassword, Server.BaseDbHost, Server.BaseDbPort, "information_schema", "parseTime=true&loc=Local&interpolateParams=true")
	return s
}
func InitDB() {
	once.Do(func() {
		BaseDB = initSqlxDB(GetBaseDSN(), "[BASE_DB] -> ", 10, 10)
		BaseInfo = initSqlxDB(GetBaseInfoDSN(), "[BASE_DB] -> ", 10, 10)
		fmt.Println("Init DB success.")
	})
}

func initSqlxDB(dbConfig, logHeader string, maxOpen, maxIdle int) *sqlx.DB {
	fmt.Printf("dbConfig: %s, logHeader: %s, maxOpen: %d, maxIdle: %d \n", dbConfig, logHeader, maxOpen, maxIdle)
	db := sqlx.MustConnect(driverName, dbConfig)
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(0)
	return db
}

func GetAllTableName() {
	//BaseDB.Query()
}

func init() {
	var filePaht = "./conf.json"

	data, err := ioutil.ReadFile(filePaht)
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}

func GetDataMap(Db *sqlx.DB, sqlstatement string) ([]map[string]interface{}, error) {
	rows, err := Db.Query(sqlstatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	return tableData, nil
}

func Getjson(Db *sqlx.DB, sqlstatement string) (string, error) {
	rows, err := Db.Query(sqlstatement)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	jsonData, err := json.MarshalIndent(tableData, "", "")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
