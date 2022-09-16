package tools

import (
	"context"
	"saavuu/config"
	"saavuu/redis"
)

func DataMigrateDemo() {
	type (
		AcceleroHeartBeat struct {
			//use client's time
			StartTime int64
			EndTime   int64
			//relative time to startTime
			AcceleroSlots []bool
			//put heartbeats here to easily return to client
			HeartBeat []uint8
			//user to predict heart rate variation
			HeartbeatPrediction []uint8
		}
		AcceleroHeartBeatNew struct {
			//use client's time
			StartTime int64
			EndTime   int64
			//relative time to startTime
			AcceleroSlots []bool
			//put heartbeats here to easily return to client
			HeartBeat []uint8
			//user to predict heart rate variation
			HeartbeatPrediction []uint16
		}
	)
	Ctx := context.Background()
	keys, _, _ := config.Cfg.Rds.Scan(Ctx, 0, "AcceleroHeartbeat:*", 1000).Result()
	his := &AcceleroHeartBeat{}
	for _, key := range keys {
		//get fields
		fields, _, err := config.Cfg.Rds.HScan(Ctx, key, 0, "*", 1024*1024).Result()
		if err != nil {
			continue
		}
		//read every field
		for _, field := range fields {
			err = redis.HGet(Ctx, config.Cfg.Rds, key, field, his)
			if err != nil {
				continue
			}
			hisNew := &AcceleroHeartBeatNew{StartTime: his.StartTime, EndTime: his.EndTime, AcceleroSlots: his.AcceleroSlots, HeartBeat: his.HeartBeat}
			// copy his to hisNew
			//convert HeartbeatPrediction
			hisNew.HeartbeatPrediction = make([]uint16, len(his.HeartbeatPrediction))
			for i, v := range his.HeartbeatPrediction {
				hisNew.HeartbeatPrediction[i] = uint16(v)
			}
			//save to new key
			newKey := "AcceleroHeartbeat1:" + key[18:]
			redis.HSet(Ctx, config.Cfg.Rds, newKey, field, hisNew)

		}
	}

}
