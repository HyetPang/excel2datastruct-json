package field

import "fmt"

type String struct {
	*TypeBase
}

// 重新实现value方法
func (i Int) Value(index int) interface{} {
	if len(i.Values) <= index || index  < 0 {
		log.Printf("列%s索引越界,返回类型%s零值\n",i.Name,i.Type)
		return 0
	}
	return fmt.Sprintf("\"%s\"",i.Values[index])
}