package api

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/tools"
)

var apiCounter tools.Counter = tools.Counter{}

func reportStates() {
	for i := 0; ApiServices.Count() == 0 && i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Second)

	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = apiServiceNames()
	log.Info().Strs(fmt.Sprintf("there are %v apis:", len(serviceNames)), serviceNames).Send()
	for {
		time.Sleep(time.Second * 60)
		serviceNames = apiServiceNames()
		for _, serviceName := range serviceNames {
			if num, _ := apiCounter.Get(serviceName); num > 0 {
				log.Info().Any("serviceName", serviceName).Any("proccessed", num).Msg("Tasks processed.")
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
