package specification

import (
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

var disAllowedServiceNamesMap = map[string]bool{
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

func ApiName(ServiceName string) string {
	//remove "api:" prefix
	if len(ServiceName) >= 4 && ServiceName[:4] == "api:" {
		ServiceName = ServiceName[4:]
	}
	//remove prefix "In" from the name
	if len(ServiceName) > 2 && strings.ToLower(ServiceName[:2]) == "in" {
		ServiceName = ServiceName[2:]
	}
	if _, ok := disAllowedServiceNamesMap[ServiceName]; ok {
		log.Panic().Msg(ServiceName + ":ServiceName misnamed. Check your code")
	}
	if len(ServiceName) == 0 {
		log.Panic().Msg("Empty ServiceName is empty")
	}
	//first byte of ServiceName should be lower case
	if ServiceName[0] >= 'A' && ServiceName[0] <= 'Z' {
		ServiceName = string(ServiceName[0]+32) + ServiceName[1:]
	}
	//ensure ServiceKey start with "api:"
	return "api:" + ServiceName
}

func TypeName(i interface{}) (name string) {
	//get default ServiceName
	var _type reflect.Type
	//take name of type v as key
	for _type = reflect.TypeOf(i); _type.Kind() == reflect.Ptr; _type = _type.Elem() {
	}
	return name

}
