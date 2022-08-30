package service

import (
	"fmt"
	"saavuu/config"
	"saavuu/http"
	. "saavuu/redis"
	"time"

	redis "github.com/go-redis/redis/v8"
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
	}
	type Output struct {
		Heartbeat uint8
	}

	http.NewService("svc:Acceleros1s", func(svcCtx *http.HttpContext) (data interface{}, err error) {
		//get jwtid
		JwtID, ok := svcCtx.JwtField("id").(string)
		if !ok {
			return nil, http.ErrJWT
		}

		var (
			in  = &Input{}
			out = &Output{}
		)
		if err = http.ToStruct(svcCtx.Req, in); err != nil || len(in.Accleration1s)%3 != 0 {
			return nil, err
		}
		now := time.Now().Unix()
		his := &AcceleroHeartBeat{StartTime: now, EndTime: now}
		if err = HGet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his); err != nil {
			return nil, err
		}

		// convert to HeartRate1s
		if err = Do(svcCtx.Ctx, config.Cfg.Rds, "svc:Acceleros1sToHeartBeat", in, out); err != nil && err != redis.Nil {
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
