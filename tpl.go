package main

const (
	LoadTpl = `
/*
    自动生成文件 请勿改动
*/
{{range $i, $v := .list}}
    var {{$v}} =require('{{$v}}');
{{end}}

{{range $i, $v := .list}}
    window.{{$v | ExportColumn}}= {{$v}};
{{end}}
`
)
