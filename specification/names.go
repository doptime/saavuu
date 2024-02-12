package specification

import (
	"fmt"
	"reflect"
	"strings"
)

var DisAllowedServiceNames = map[string]bool{
	"string":     true,
	"int32":      true,
	"int64":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"uint":       true,
	"float":      true,
	"bool":       true,
	"byte":       true,
	"rune":       true,
	"complex64":  true,
	"complex128": true,
}

// 不能Panic,因为可能被web客户端调用。否则服务端可以被客户端恶意关闭
func ApiName(ServiceName string) (string, error) {
	//remove "api:" prefix
	if len(ServiceName) >= 4 && ServiceName[:4] == "api:" {
		ServiceName = ServiceName[4:]
	}
	//remove prefix "In" from the name
	if len(ServiceName) > 2 && strings.ToLower(ServiceName[:2]) == "in" {
		ServiceName = ServiceName[2:]
	}
	if _, ok := DisAllowedServiceNames[ServiceName]; ok {
		return "", fmt.Errorf("service misnamed. %s", ServiceName)
	}
	if len(ServiceName) == 0 {
		return "", fmt.Errorf("service misnamed. %s", ServiceName)
	}
	//first byte of ServiceName should be lower case
	if ServiceName[0] >= 'A' && ServiceName[0] <= 'Z' {
		ServiceName = string(ServiceName[0]+32) + ServiceName[1:]
	}
	//ensure ServiceKey start with "api:"
	return "api:" + ServiceName, nil
}

func TypeName(i interface{}) (name string) {
	//get default ServiceName
	var _type reflect.Type
	//take name of type v as key
	for _type = reflect.TypeOf(i); _type.Kind() == reflect.Ptr; _type = _type.Elem() {
	}
	return name

}
