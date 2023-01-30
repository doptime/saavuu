package saavuu

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/yangkequn/saavuu/config"
)

type ServiceInfo struct {
	// ServiceName is the name of the service
	ServiceName string
	// ServiceFunc is the function of the service
	ServiceFunc func(s []byte) (err error)
	// ServiceBatchSize is the number of data to be fetched from redis at one time
	ServiceBatchSize int64
}

var services map[string]*ServiceInfo = map[string]*ServiceInfo{}

func pipingServiceTask() {
	var (
		data         []string
		taskReceived int
		ServiceInfos []*ServiceInfo = make([]*ServiceInfo, 0, len(services))
		serviceInfo  *ServiceInfo
	)
	c := context.Background()
	//from services to ServiceInfos
	for _, serviceInfo = range services {
		ServiceInfos = append(ServiceInfos, serviceInfo)
	}

	for {
		//fetch datas from redis,using LRange
		pipline := config.ParamRds.Pipeline()
		for _, serverInfo := range ServiceInfos {
			pipline.LRange(c, serverInfo.ServiceName, 0, serverInfo.ServiceBatchSize-1)
			pipline.LTrim(c, serverInfo.ServiceName, serverInfo.ServiceBatchSize, -1)
			ServiceInfos = append(ServiceInfos, serverInfo)
		}
		cmd, err := pipline.Exec(c)

		for i, cmd := range cmd {
			//skip LTrim
			if i%2 == 1 {
				continue
			}
			//skip LRange if no data
			if data = cmd.(*redis.StringSliceCmd).Val(); err != nil || len(data) == 0 {
				continue
			}
			//nolonger using BLPop to receive another 1 data, avoid sockert timeout as service increase
			serviceInfo = ServiceInfos[i/2]
			for _, s := range data {
				go serviceInfo.ServiceFunc([]byte(s))
				serviceCounter.Add(serviceInfo.ServiceName, 1)
				taskReceived++
			}
		}
		if taskReceived == 0 {
			time.Sleep(time.Millisecond * 32)
		} else {
			taskReceived = 0
		}
	}
}
