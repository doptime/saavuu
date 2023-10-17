package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/tools"
)

var apiCounter tools.Counter = tools.Counter{}

func reportStates() {
	time.Sleep(time.Second * 5)
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = apiServiceNames()
	log.Info().Strs(fmt.Sprintf("there are %v apis:", len(serviceNames)), serviceNames).Send()
	for {
		time.Sleep(time.Second * 60)
		now := time.Now().String()[11:19]
		for _, serviceName := range serviceNames {
			if num, _ := apiCounter.Get(serviceName); num > 0 {
				log.Info().Msg(now + "" + serviceName + " proccessed " + strconv.Itoa(int(num)) + " tasks")
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
