package api

import (
	"strconv"
	"time"

	"github.com/yangkequn/saavuu/logger"
	"github.com/yangkequn/saavuu/tools"
)

var apiCounter tools.Counter = tools.Counter{}

func reportStates() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(apiServices))
	for serviceName := range apiServices {
		serviceNames = append(serviceNames, serviceName)
	}
	logger.Lshortfile.Println(len(serviceNames), "apis:", serviceNames)
	for {
		time.Sleep(time.Second * 60)
		now := time.Now().String()[11:19]
		for _, serviceName := range serviceNames {
			if num, _ := apiCounter.Get(serviceName); num > 0 {
				logger.Lshortfile.Println(now + "" + serviceName + " proccessed " + strconv.Itoa(int(num)) + " tasks")
			}
			apiCounter.DeleteAndGetLastValue(serviceName)
		}
	}
}
func RunningAllApis() {
	delayTasksLoad()
	go reportStates()
	receiveJobs()
}
