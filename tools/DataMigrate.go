package tools

import (
	"context"
	"saavuu/config"
	"strconv"
)

func DataMigrateDemo() (err error) {
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
			HeartbeatPrediction []uint16
		}
		AccHrTraj struct {
			//use client's time
			Start  int64
			End    int64
			HrCnt  int32
			AccCnt int32
		}
	)
	Ctx := context.Background()
	keys, _, _ := config.Cfg.ParamRedis.Scan(Ctx, 0, "AcceleroHeartbeat:*", 1000).Result()
	his := &AcceleroHeartBeat{}
	for _, key := range keys {
		//get fields
		fields, _, err := config.Cfg.ParamRedis.HScan(Ctx, key, 0, "*", 1024*1024).Result()
		if err != nil {
			continue
		}
		//read every field
		for _, field := range fields {
			err = HGet(Ctx, config.Cfg.ParamRedis, key, field, his)
			if err != nil {
				continue
			}
			var StartTime string = strconv.FormatInt(his.StartTime, 10)
			hisNew := &AccHrTraj{Start: his.StartTime, End: his.EndTime}
			//count HBCnt
			hisNew.HrCnt = 0
			for _, hb := range his.HeartBeat {
				if hb > 0 {
					hisNew.HrCnt++
				}
			}
			//count AccCnt
			hisNew.AccCnt = 0
			for _, acc := range his.AcceleroSlots {
				if acc {
					hisNew.AccCnt++
				}
			}
			//save to new key
			// newKey := "TrajAccHr:" + key[18:]
			// redis.HSet(Ctx, config.Cfg.Rds, newKey, field, hisNew)

			//move HeartBeat to new key
			HeartBeatKey := "TrajHr:" + key[18:]
			HSet(Ctx, config.Cfg.ParamRedis, HeartBeatKey, StartTime, his.HeartBeat)

			//move AcceleroSlots to new key
			// AcceleroSlotsKey := "TrajAcc:" + key[18:] + ":" + StartTime
			// for i, acc := range his.AcceleroSlots {
			// 	ok := false
			// 	if acc {
			// 		FormerAcceleroKey := "Accelero:" + key[18:]
			// 		acceleroData := &[]int16{}
			// 		var timeSlot string = strconv.FormatInt(his.StartTime+int64(i), 10)
			// 		if err = redis.HGet(Ctx, config.Cfg.Rds, FormerAcceleroKey, timeSlot, acceleroData); err == nil {
			// 			redis.RPush(Ctx, config.Cfg.Rds, AcceleroSlotsKey, acceleroData)
			// 			ok = true
			// 		}
			// 	}
			// 	if !ok {
			// 		redis.RPush(Ctx, config.Cfg.Rds, AcceleroSlotsKey, nil)
			// 	}
			// }
		}
	}
	return nil
}
