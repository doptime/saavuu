package api

import (
	"fmt"

	"github.com/yangkequn/saavuu/specification"
)

func CallByHTTP(ServiceName string, paramIn map[string]interface{}) (ret interface{}, err error) {
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
		var rpc = Rpc[interface{}, interface{}](OptName(ServiceName), OptDb(apiInfo.DbName))
		return rpc(paramIn)
	}
	//if function is stored locally, call it directly. This is alias monolithic mode
	if buf, err = specification.MarshalApiInput(paramIn); err != nil {
		return nil, err
	}
	return apiInfo.ApiFuncWithMsgpackedParam(buf)
}
