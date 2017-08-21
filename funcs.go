package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FirstCharUpper(str string) string {
	if len(str) > 0 {
		return strings.ToUpper(str[0:1]) + str[1:]
	} else {
		fmt.Println("111111111")
		return ""
	}
}

func ExportColumn(columnName string) string {
	columnItems := strings.Split(columnName, "_")
	columnItems[0] = FirstCharUpper(columnItems[0])
	for i := 0; i < len(columnItems); i++ {
		item := strings.Title(columnItems[i])

		if strings.ToUpper(item) == "ID" {
			item = "ID"
		}

		columnItems[i] = item
	}

	return strings.Join(columnItems, "")
}

func AddInt(i int, b bool) int {
	var t = 1
	if b {
		t = 1
	} else {
		t++
	}

	return t

}

func TypeConvert(str string, value interface{}) interface{} {
	switch str {
	case "smallint", "tinyint", "int", "bigint", "float", "double", "decimal":

		switch value.(type) {
		case string:
			v, err := strconv.Atoi(value.(string))
			if err != nil {
				panic(err)
			}
			return v
		default:
			return value.(int)
		}
	case "timestamp", "datetime":
		v := value.(time.Time)
		return v.Format("2006-01-02 15:04:05")
	default:
		switch value.(type) {
		case int:
			return strconv.Itoa(value.(int))
		default:
			return value.(string)
		}
	}
}

///
