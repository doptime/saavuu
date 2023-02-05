package api

import (
	"strconv"
	"time"

	"github.com/yangkequn/saavuu/logger"
	"github.com/yangkequn/saavuu/tools"
)

var apiCounter tools.Counter = tools.Counter{}

func reportApiStates() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(apiServices))
	for serviceName := range apiServices {
		serviceNames = append(serviceNames, serviceName)
	}
	logger.Lshortfile.Println("service has", len(serviceNames), "services:", serviceNames)
	for {
		time.Sleep(time.Second * 60)
		now := time.Now().String()[11:19]
		for _, serviceName := range serviceNames {
			num, _ := apiCounter.Get(serviceName)
			logger.Lshortfile.Println(now + " service " + serviceName + " proccessed " + strconv.Itoa(int(num)) + " tasks")
			apiCounter.DeleteAndGetLastValue(serviceName)
		}
	}
}
func RunningAllApis() {
	PendingApiFromRedisToLoal()
	go reportApiStates()
	receiveApiJobs()
}
