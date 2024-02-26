package api

import (
	"context"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
)

type ApiInfo struct {
	// Name is the name of the service
	Name       string
	DataSource string
	Ctx        context.Context
	// ApiFuncWithMsgpackedParam is the function of the service
	ApiFuncWithMsgpackedParam func(s []byte) (ret interface{}, err error)
}

var ApiServices cmap.ConcurrentMap[string, *ApiInfo] = cmap.New[*ApiInfo]()

func apiServiceNames() (serviceNames []string) {
	for _, serviceInfo := range ApiServices.Items() {
		serviceNames = append(serviceNames, serviceInfo.Name)
	}
	return serviceNames
}
func GetServiceDB(serviceName string) *redis.Client {
	serviceInfo, _ := ApiServices.Get(serviceName)
	DataSourceName := serviceInfo.DataSource
	return config.Rds[DataSourceName]
}

var fun2ApiInfoMap = &sync.Map{}
var APIGroupByDataSourceName = cmap.New[[]string]()
