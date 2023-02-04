package saavuu

import (
	"strconv"
	"time"

	"github.com/yangkequn/saavuu/logger"
)

var rpcCounter Counter = Counter{}

func reportRpcStates() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(rpcServices))
	for serviceName := range rpcServices {
		serviceNames = append(serviceNames, serviceName)
	}
	logger.Lshortfile.Println("service has", len(serviceNames), "services:", serviceNames)
	for {
		time.Sleep(time.Second * 60)
		now := time.Now().String()[11:19]
		for _, serviceName := range serviceNames {
			num, _ := rpcCounter.Get(serviceName)
			logger.Lshortfile.Println(now + " service " + serviceName + " proccessed " + strconv.Itoa(int(num)) + " tasks")
			rpcCounter.DeleteAndGetLastValue(serviceName)
		}
	}
}
func RunningAllRpcs() {
	PendingRpcFromRedisToLoal()
	go reportRpcStates()
	receiveRpcJobs()
}
