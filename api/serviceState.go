package api

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/tools"
)

var apiCounter tools.Counter = tools.Counter{}

func reportStates() {
	//wait till all apis are loaded
	for i, lastCnt := 0, 0; (ApiServices.Count() == 0 || lastCnt != ApiServices.Count()) && i < 100; i++ {
		lastCnt = ApiServices.Count()
		time.Sleep(time.Millisecond * 100)
	}

	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = apiServiceNames()
	log.Info().Any("cnt", len(serviceNames)).Strs("apis are load:", serviceNames).Send()
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
