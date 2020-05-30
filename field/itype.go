package field

import "fmt"

// excel定义的字段
type Field interface {
	// 类型字符串
	// Type() string
	// 获取注释
	// Comment() string
	// 字段名
	GetName() string
	// 字段值
	Value(index int) interface{}
	// 设置值
	AddValue(value string)
	// 获取值长度
	ValueLen() int
}

// 类型基类
type TypeBase struct {
	// 注释
	Comment string
	// 字段名字
	Name string
	// 类型
	Type string
	// 类型对应的值,整个列的值
	Values []string
}

// 类型字符串
// func (tb TypeBase) Type() string {
// 	return tb.Type
// }

// // 获取注释
// func (tb TypeBase) Comment() string {
// 	return tb.Comment
// }

// 字段名
func (tb TypeBase) GetName() string {
	return fmt.Sprintf("\"%s\"",tb.Name)
}

// 字段值
// func (tb TypeBase) Value(index int) interface{} {
	// return tb.Values[index]
// }

// 字段值
func (tb *TypeBase) AddValue(value string) {
	tb.Values = append(tb.Values,value)
}

// 字段值
func (tb TypeBase) ValueLen() int {
	return len(tb.Values)
}