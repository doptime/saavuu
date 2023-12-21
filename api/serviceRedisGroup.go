package api

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
)

func defaultXReadGroupArgs() *redis.XReadGroupArgs {
	var (
		streams []string
	)
	services := apiServiceNames()
	streams = append(streams, services...)
	//from services to ServiceInfos
	for i := 0; i < len(services); i++ {
		//append default stream id
		streams = append(streams, ">")
	}

	//ServiceBatchSize is the number of tasks that a service can read from redis at the same time
	args := &redis.XReadGroupArgs{Streams: streams, Block: time.Second * 20, Count: config.Cfg.Api.ServiceBatchSize, NoAck: true, Group: "group0", Consumer: "saavuu"}
	return args
}
func XGroupCreateOne(c context.Context, serviceName string) (err error) {
	var (
		rds *redis.Client = config.RdsClientDefault()
	)

	//continue if the group already exists
	if cmd := rds.XInfoGroups(c, serviceName); cmd.Err() == nil || len(cmd.Val()) > 0 {
		return nil
	}
	//create a group if none exists
	if cmd := rds.XGroupCreateMkStream(c, serviceName, "group0", "$"); cmd.Err() != nil {
		log.Info().AnErr("XGroupCreateOne", cmd.Err()).Send()
		return cmd.Err()
	}
	return nil
}
