package specification

import (
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

var DisAllowedServiceNames = map[string]bool{
	"":           true,
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

// return the api name of the service
// name with format "api:serviceName". first letter of serviceName should be lower case, and start with "api:"
// two possible source of the service name:
// 1. the type name of the first parameter of the function
// 2. the name give by the user
// do not panic, because it may be called by web client. otherwise the server can be maliciously closed by the client
func ApiName(ServiceName string) string {
	//remove  prefix. "api:" is the case of encoded service name. other wise for the case of parameter type name
	var prefixes = []string{"api:", "input", "in", "req", "arg", "param", "src", "data"}
	if ServiceNameLowercase := strings.ToLower(ServiceName); len(ServiceNameLowercase) > 0 {
		for _, prefix := range prefixes {
			if strings.HasPrefix(ServiceNameLowercase, prefix) {
				ServiceName = ServiceName[len(prefix):]
			}
		}
	}
	if _, ok := DisAllowedServiceNames[ServiceName]; ok {
		log.Error().Str("service misnamed", ServiceName).Send()
		return ""
	}
	//first byte of ServiceName should be lower case
	if ServiceName[0] >= 'A' && ServiceName[0] <= 'Z' {
		ServiceName = string(ServiceName[0]+32) + ServiceName[1:]
	}
	//ensure ServiceKey start with "api:"
	return "api:" + ServiceName
}

func ApiNameByType(i interface{}) (name string) {
	//get default ServiceName
	var _type reflect.Type
	//take name of type v as key
	for _type = reflect.TypeOf(i); _type.Kind() == reflect.Ptr; _type = _type.Elem() {
	}
	return ApiName(_type.Name())

}
