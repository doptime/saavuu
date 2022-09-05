package service

import (
	"fmt"
	"saavuu/https"
)

func PrintServices() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(https.ServiceMap))
	for k := range https.ServiceMap {
		serviceNames = append(serviceNames, k)
	}

	fmt.Println("ServiceMap has", len(https.ServiceMap), "services:", serviceNames)
}
