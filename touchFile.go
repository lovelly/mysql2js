package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strings"
)

func TouchDir() {
	cmd := exec.Command("cmd", "/C", "rd", "/s", "/q", "template")
	cmd.Dir = "./"
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("remove template not foud foldor :", err.Error())
	}
	fmt.Println(string(out))

	cmd = exec.Command("cmd", "/C", "mkdir", "template")
	cmd.Dir = "./"
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Println("RanderDB error :", err.Error())
		return
	}

	fmt.Println(string(out))
}

func Rander(tableName string, t *template.Template) {
	tableInfo, terr := GetDataMap(BaseInfo, fmt.Sprintf(tableInfoSql, tableName))
	if terr != nil {
		panic(terr)
	}
	data, derr := GetDataMap(BaseDB, fmt.Sprintf("select * from %s;", tableName))
	if derr != nil {
		panic(derr)
	}

	fileName := "template/" + strings.ToLower(tableName) + ".js"
	os.Remove(fileName)
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	var Ments []*Comment
	var Keys []string
	var ColType = make(map[string]string)
	for _, v := range tableInfo {
		m := &Comment{}
		COLUMN_NAME := v["COLUMN_NAME"].(string)
		m.Field = COLUMN_NAME
		m.Ment = v["COLUMN_COMMENT"].(string)
		Ments = append(Ments, m)
		if v["COLUMN_KEY"].(string) == "PRI" {
			Keys = append(Keys, ExportColumn(COLUMN_NAME))
		}
		ColType[COLUMN_NAME] = v["DATA_TYPE"].(string)
	}

	s := ReSetData(Keys, data, ColType)
	s = Comments(Ments) + "var " + tableName + " = " + s + QueruFunc(tableName)
	write := bufio.NewWriter(f)
	write.WriteString(s)
	write.Flush()
	//if err := t.Execute(f, map[string]interface{}{
	//	"data":      data,
	//	"info":      Ments,
	//	"TableName": tableName,
	//}); err != nil {
	//	panic(err)
	//}
}

func ReSetData(Keys []string, data []map[string]interface{}, ColType map[string]string) string {
	//var t = map[interface{}]interface{}{}
	if len(data) < 1 {
		return "[]"
	}
	for _, item := range data {
		for k, v := range item {
			newK := ExportColumn(k)
			delete(item, k)
			item[newK] = TypeConvert(ColType[k], v)
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

func Comments(Ments []*Comment) string {
	s := "/*\n	自动生成文件 请勿改动 \n"
	for _, v := range Ments {
		s = s + v.Field + " : " + v.Ment + "\n"
	}

	s += "*/\n"
	return s
}

func QueruFunc(tableName string) string {
	s := "\n" + tableName + `.prototype.Query = function (param) {
	var ret = []
	for (item in %s) {
		var has = true;
		for (k in param) {
			if (item[k] != param[k]){
				has = false;
				break;
			}
		}
		if (has) {
			ret.push(item)
		}
	}
	return ret
}
` +
		tableName + `.prototype.QueryOne = function (param) {
	for (item in %s) {
		var has = true;
		for (k in param) {
			if (item[k] != param[k]){
				has = false;
				break;
			}
		}
		if (has) {
			return item
		}
	}
}
` + `module.exports = ` + tableName
	return fmt.Sprintf(s, tableName, tableName)
}


func RanderLoad(list []string) {
	fileName := "template/" + strings.ToLower("loadTmplate") + ".js"
	os.Remove(fileName)
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("load_template").
		Funcs(template.FuncMap{
			"FirstCharUpper": FirstCharUpper,
			"ExportColumn":   ExportColumn,
			"AddInt":         AddInt,
		}).
		Parse(LoadTpl))
	t.Execute(f, map[string]interface{}{
		"list": list,
	})
}

func RanderTemplate(list []string) {
	t := template.Must(template.New("template").
		Funcs(template.FuncMap{
			"FirstCharUpper": FirstCharUpper,
			"ExportColumn":   ExportColumn,
			"AddInt":         AddInt,
		}).
		Parse(LoadTpl))

	for _, v := range list {
		Rander(v, t)
	}
}
