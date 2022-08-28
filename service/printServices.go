package service

import (
	"fmt"
	"saavuu/http"
)

func PrintServices() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(http.ServiceMap))
	for k := range http.ServiceMap {
		serviceNames = append(serviceNames, k)
	}

	fmt.Println("ServiceMap has", len(http.ServiceMap), "services:", serviceNames)
}
