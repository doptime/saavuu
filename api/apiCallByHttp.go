package api

import (
	"fmt"

	"github.com/yangkequn/saavuu/specification"
)

func CallByHTTP(ServiceName string, paramIn map[string]interface{}) (ret interface{}, err error) {
	var (
		fuc *ApiInfo
		ok  bool
		buf []byte
	)
	if ServiceName = specification.ApiName(ServiceName); len(ServiceName) == 0 {
		return nil, fmt.Errorf("service misnamed %s", ServiceName)
	}
	var rpc = Rpc[interface{}, interface{}](OptName(ServiceName))
	//if function is stored locally, call it directly. This is alias monolithic mode
	if fuc, ok = ApiServices.Get(ServiceName); ok {
		if buf, err = specification.MarshalApiInput(paramIn); err != nil {
			return nil, err
		}
		return fuc.ApiFuncWithMsgpackedParam(buf)
	} else {
		//if function is not stored locally, call it remotely (RPC). This is alias microservice mode
		return rpc(paramIn)
	}

}