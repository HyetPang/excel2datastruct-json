package field

import (
	"log"
)

// 返回对应的类型处理器
func GetTypeHandler(typeString,comment,name string) Field {
	switch typeString {
	case "int","uint","int32","uint32","int64","uint64":
		return Int{
			TypeBase:&TypeBase{
				Type:typeString,
				Comment:comment,
				Name:name,
				Values:make([]string,0,50),
			},
		}
	case "float32","float64":
		return Float{
			TypeBase:&TypeBase{
			Type:typeString,
			Comment:comment,
			Name:name,
			Values:make([]string,0,50),
		},}
	case "string":
		return String{
			TypeBase:&TypeBase{
				Type:typeString,
				Comment:comment,
				Name:name,
				Values:make([]string,0,50),
			},
		}
	}
	log.Printf("未知的类型:%s",typeString)
	return nil
}