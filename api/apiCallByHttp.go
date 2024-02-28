package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/yangkequn/saavuu/aopt"
	"github.com/yangkequn/saavuu/specification"
)

func CallByHTTP(ServiceName string, paramIn map[string]interface{}, req *http.Request) (ret interface{}, err error) {
	var (
		apiInfo *ApiInfo
		ok      bool
		buf     []byte
	)
	if ServiceName = specification.ApiName(ServiceName); len(ServiceName) == 0 {
		return nil, fmt.Errorf("service misnamed %s", ServiceName)
	}
	//if function is stored locally, call it directly. This is alias monolithic mode
	if apiInfo, ok = ApiServices.Get(ServiceName); !ok {
		//if function is not stored locally, call it remotely (RPC). This is alias microservice mode
		var rpc = Rpc[interface{}, interface{}](aopt.Name(ServiceName), aopt.DataSource(apiInfo.DataSource))
		return rpc(paramIn)
	}
	if apiInfo.WithHeader {
		//copy fields from req to paramIn
		for key, value := range req.Header {
			if len(value) > 1 {
				paramIn["Header"+key] = value
			} else {
				paramIn["Header"+key] = value[0]
			}
		}
		// copy ip address from req to paramIn
		paramIn["Header"+"Ip"] = req.RemoteAddr
		paramIn["Header"+"Host"] = req.Host
		paramIn["Header"+"Method"] = req.Method
		paramIn["Header"+"Path"] = req.URL.Path
		paramIn["Header"+"Query"] = req.URL.RawQuery

	}
	//if function is stored locally, call it directly. This is alias monolithic mode
	if buf, err = specification.MarshalApiInput(paramIn); err != nil {
		return nil, err
	}
	return apiInfo.ApiFuncWithMsgpackedParam(buf)
}

func HeaderFieldsUsed[i any](param i) bool {
	var (
		vType reflect.Type
	)
	//use reflect to detect if the param has a field start with "Header", or tag of that field contains "Header",if true return true else return false

	// case double pointer decoding
	if vType = reflect.TypeOf((*i)(nil)).Elem(); vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}
	for i := 0; i < vType.NumField(); i++ {
		name := strings.ToLower(vType.Field(i).Name)
		tag := strings.ToLower(vType.Field(i).Tag.Get("mapstructure"))
		if strings.HasPrefix(name, "header") || strings.Contains(tag, "header") {
			return true
		}
	}
	return false
}
