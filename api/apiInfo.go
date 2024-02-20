package api

import (
	"context"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type ApiInfo struct {
	// ApiName is the name of the service
	ApiName string
	DbName  string
	Ctx     context.Context
	// ApiFuncWithMsgpackedParam is the function of the service
	ApiFuncWithMsgpackedParam func(s []byte) (ret interface{}, err error)
}

var ApiServices cmap.ConcurrentMap[string, *ApiInfo] = cmap.New[*ApiInfo]()

func apiServiceNames() (serviceNames []string) {
	for _, serviceInfo := range ApiServices.Items() {
		serviceNames = append(serviceNames, serviceInfo.ApiName)
	}
	return serviceNames
}

var fun2ApiInfoMap = &sync.Map{}
