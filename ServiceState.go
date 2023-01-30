package saavuu

import (
	"strconv"
	"time"

	"github.com/yangkequn/saavuu/logger"
)

var serviceCounter Counter = Counter{}

func reportServiceStates() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(services))
	for serviceName := range services {
		serviceNames = append(serviceNames, serviceName)
	}
	logger.Lshortfile.Println("service has", len(serviceNames), "services:", serviceNames)
	for {
		time.Sleep(time.Second * 60)
		now := time.Now().String()[11:19]
		for _, serviceName := range serviceNames {
			num, _ := serviceCounter.Get(serviceName)
			logger.Lshortfile.Println(now + " service " + serviceName + " proccessed " + strconv.Itoa(int(num)) + " tasks")
			serviceCounter.DeleteAndGetLastValue(serviceName)
		}
	}
}
func RunningAllService() {
	go reportServiceStates()
	pipingServiceTask()
}
