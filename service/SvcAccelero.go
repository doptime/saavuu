package service

import (
	"fmt"
	"saavuu/config"
	"saavuu/http"
	"saavuu/redis"
)

func init() {
	type AcceleroHeartBeat struct {
		StartTime     int64
		EndTime       int64
		Accleration1s []int16
		HeartBeat     []uint8
		HBPrediction  []uint8
	}
	type Input struct {
		Accleration1s []int16
		HeartBeat     uint8
	}
	type Output struct {
		Heartbeat uint8
	}

	http.NewService("svc:Acceleros1s", func(svcCtx *http.HttpContext) (data interface{}, err error) {
		var (
			in  = &Input{}
			out = &Output{}
		)
		if err = http.ToStruct(svcCtx.Req, in); err != nil || len(in.Accleration1s)%3 != 0 {
			return nil, err
		}
		// convert to HeartRate1s
		if err = redis.Do(svcCtx.Ctx, config.Cfg.Rds, "svc:Acceleros1sToHeartBeat", in, out); err != nil {
			return nil, err
		}
		fmt.Println("lengthof data is ", len(in.Accleration1s)/3)
		for i := 0; i < len(in.Accleration1s)/3; i += 3 {
			x := float32(in.Accleration1s[i]) * 16.0 / 32768.0
			y := float32(in.Accleration1s[i+1]) * 16.0 / 32768.0
			z := float32(in.Accleration1s[i+2]) * 16.0 / 32768.0
			fmt.Println("x is ", x, "y is ", y, "z is ", z)
		}
		return data, nil
	})
}
