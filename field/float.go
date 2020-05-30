package field

import (
	"log"
	"strconv"
)

type Float struct {
	*TypeBase
}

// 重新实现value方法
func (i Int) Value(index int) interface{} {
	if len(i.Values) <= index || index  < 0 {
		log.Printf("列%s索引越界,返回类型%s零值\n",i.Name,i.Type)
		return 0
	}
	floatValue,err := strconv.ParseFloat(i.Values[index],64)
	if err != nil {
		log.Printf("列%的第%d行的值不是数字,类型转换出错:%s,返回0\n",i.Name,index+3+1,err.Error())
		return 0
	}
	return floatValue
}