package main

import (
	"excel2datastructjson/field"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/tealeg/xlsx"
)

var (
	filepath   = flag.String("file", "", "excel文件的绝对路径")
	format   = flag.String("format", "array", "json数据的格式类型,支持两种,dict和array,默认是array")
)

// 用法函数
func useage() {
	fmt.Fprintf(os.Stderr, "Usage of excel2json:\n")
	fmt.Fprintf(os.Stderr, "\texcel2json [flags] -file excel文件\n")
	fmt.Fprintf(os.Stderr, "\texcel2json [flags] -format [dict,array]\n")
	fmt.Fprintf(os.Stderr, "\t只接受以.xlsx后缀的excel文件,并且文件的前三行以后才是数据,前三行分别是字段注释,字段类型,字段名字,字段类型,目前支持uint,uint32,uint64,int,int32,int64,float32,float64,string\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://github.com/HyetPang/excel2datastruct-json\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	windowEnter = "\r\n"
	linuxEnter = "\n"
	maxEnter = "\n"
	formatArrayStart = "["
	formatArrayEnd = "]"
	formatDictStart = "{"
	formatDictEnd = "}"
	formatArray = "array"
	formatDict = "dict"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("excel2json: ")
	flag.Usage = useage
	flag.Parse()
	if len(*filepath) == 0 {
		// 没有给定文件路径
		flag.Usage()
		os.Exit(2)
	}
	if format != formatArray && format != formatDict {
		flag.Usage()
		os.Exit(2)
	}
	if !strings.HasSuffix(*filepath,".xlsx") {
		// 只接受xlsx的文件
		flag.Usage()
		os.Exit(2)
	}
	handle()
}

func handle() {
	// 打开文件
	file,err := xlsx.OpenFile(filepath)
	if err != nil {
		log.Printf("文件%s打开出错:%s",filepath,err.Error())
		os.Exit(2)
	}
	for _, sheet := range file.Sheets {
		// 处理数据
		handleData(sheet)
	}
}

func getData(sheet *xlsx.Sheet) {
	// 处理前三行元数据
	rows := sheet.Rows
	if len(rows) <= 3 {
		log.Printf("工作簿%s数据行数应该大于3行,忽略这个工作簿数据的处理",sheet.Name)
		return
	}
	fields := make([]*field.Field,0,10)
	// 处理数据
	data := rows[3:]
	for cellIndex, cell := range rows[1].Cells {
		comment := strings.ReplaceAll(rows[0].Cells[cellIndex].Value,"\r\n",",")
		name := rows[2].Cells[cellIndex].Value,"\r\n",","
		field := field.GetTypeHandler(cell.Value,comment,name)
		if field == nil {
			return
		}
		// 数据
		for rowIndex := range data {
			field.AddValue(sheet.Cell(rowIndex,cellIndex))
		}
		fields = append(fields,field)
	}
	// 保存数据结构
	var dataStruct strings.Builder
	dataStruct.WriteString(fmt.Sprintf("type %s struct {",upperFirst(getCamelString(sheet.Name))))
	for _, field := range fields {
		dataStruct.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`",upperFirst(getCamelString(field.Name())),field.Type(),field.Name()))
	}
	dataStruct.WriteString("}")
	err := ioutil.WriteFile(fmt.Sprintf("%s.go",sheet.Name),[]byte(dataStruct.String()),0755)
	if err != nil {
		log.Printf("工作簿%s生成go文件出错,err:%s\n",sheet.Name,err.Error())
		return
	}
	// 保存json
	var json strings.Builder
	if format == formatArray {
		json.WriteString(formatArrayStart)
	} else {
		json.WriteString(formatDictStart)
	}
	dataLen := len(fields[0].ValueLen())
	for fieldIndex, field := range fields {
		if format == formatDict {
			json.WriteString("\t\"%s\": ",field.Value(i))
		}
		json.WriteString(fmt.Sprintf("%s%s",formatDictStart,linuxEnter))
		for i:=0;i<dataLen;i++ {
			json.WriteString(fmt.Sprintf("\t\t%s: ",field.Name()))
			json.WriteString(field.Value(i))
			if i+1 != dataLen {
				json.WriteString(",")
			}
			json.WriteString(fmt.Sprintf("%s",linuxEnter))
		}
		json.WriteString(fmt.Sprintf("%s%s",formatDictEnd,linuxEnter))
		json.WriteString(fmt.Sprintf("%s",formatDictEnd))
		if fieldIndex+1 != len(fields) {
			json.WriteString(",")
		}
		json.WriteString(fmt.Sprintf("%s",linuxEnter))
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s.json",sheet.Name),[]byte(json.String()),0755)
	if err != nil {
		log.Printf("工作簿%s生成json文件出错,err:%s\n",sheet.Name,err.Error())
		return
	}
}

func getCamelString(ss string) string {
	temp := strings.Split(ss, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') {
				vv[0] -= 32
			}
			s += string(vv)
		}
	}
	return s
}

func upperFirst(s string) string {
	runes := []rune(s)
	if bool(runes[0] >= 'a' && runes[0] <= 'z') {
		runes[0] -= 32
	}
	string(runes)
}
