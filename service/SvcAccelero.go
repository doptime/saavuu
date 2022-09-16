package service

import (
	"context"
	"fmt"
	"saavuu/config"
	"saavuu/https"
	. "saavuu/redis"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func init() {
	//each AcceleroHeartBeat should hold data for very long time, so allow  response to client conviniently

	//client upload Accleration1s or HeartBeat,But not both
	type (
		Input struct {
			Accleration1s []int16
			HeartRate     uint8
			//use client's time
			Time int64
		}
		AcceleroHeartBeat struct {
			//use client's time
			StartTime int64
			EndTime   int64
			//relative time to startTime
			AcceleroSlots []bool
			//put heartbeats here to easily return to client
			HeartBeat []uint8
			//heart rate variation predicted by machine learning
			//value /= 256 gets the real value of heart rate variation
			HeartbeatPrediction []uint16
		}
	)
	checkInput := func(v interface{}) error {
		in := v.(*Input)
		//validate input
		if in.Accleration1s == nil && in.HeartRate == 0 {
			return fmt.Errorf("ErrInvalidInput")
		}
		if in.Time <= 0 || absInt64(in.Time-time.Now().Unix()) > 3600*4 {
			return fmt.Errorf("ErrInvalidTimeInput")
		}
		if in.Accleration1s != nil && len(in.Accleration1s)%3 != 0 {
			return fmt.Errorf("ErrInvalidAccleration1s")
		}
		return nil
	}

	PredicHeartRate := func(ctx context.Context, cursor int64, JwtID string, in *Input) (err error) {
		// convert to HeartRate1s
		type outType struct {
			Heartbeat float32
		}
		out := &outType{}
		paramIn := map[string]interface{}{"UID": JwtID, "Time": in.Time, "Accelero1s": in.Accleration1s}
		if err = Do(ctx, config.Cfg.Rds, "svc:PredictHeartBeat", paramIn, out); err != nil {
			return err
		}
		his := &AcceleroHeartBeat{}
		if err = HGet(ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his); err != nil && err != redis.Nil {
			return err
		}
		//not allow to change history
		if his.EndTime-his.StartTime < cursor {
			return
		}

		//write predicted HeartBeat to his
		his.HeartbeatPrediction[cursor] = uint16(out.Heartbeat * 256.0)
		HSet(ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his)
		fmt.Println("PredictHeartBeat", out.Heartbeat)
		return nil
	}
	https.NewService("svc:HeartbeatAcceleros1s", func(svcCtx *https.HttpContext) (data interface{}, err error) {
		//get jwtid
		JwtID, ok := svcCtx.JwtField("id").(string)
		if !ok || len(JwtID) == 0 {
			return nil, https.ErrJWT
		}

		var (
			in = &Input{}
		)
		if err = https.ToValidStruct(svcCtx.Req, in, checkInput); err != nil {
			return nil, err
		}

		//get previous data
		his := &AcceleroHeartBeat{StartTime: in.Time, EndTime: in.Time}
		if err = HGet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his); err != nil && err != redis.Nil {
			return nil, err
		}
		// save discontinual data to history
		dataDiscontinual := in.Time-his.EndTime > 5*60
		dataTooLong := (his.EndTime - his.StartTime) > 24*3600
		if dataDiscontinual || dataTooLong {
			//save data to new key,if data less than 1 minute long,discard it
			HSet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, strconv.FormatInt(his.StartTime, 10), his)
			his.StartTime = in.Time
			his.EndTime = in.Time
			his.AcceleroSlots = make([]bool, 64)
			his.HeartBeat = make([]uint8, 64)
			//64 seconds should be enough to predict heart rate
			his.HeartbeatPrediction = make([]uint16, 64)
		}

		//make sure capacity of his.HeartBeat is enough, append 64 0s if not
		if in.Time > his.EndTime {
			his.EndTime = in.Time
		}
		if l := his.EndTime - his.StartTime + 1; int64(len(his.HeartBeat)) < l {
			his.HeartBeat = append(his.HeartBeat, make([]uint8, l-int64(len(his.HeartBeat))+64)...)
		}
		if l := his.EndTime - his.StartTime + 1; int64(len(his.AcceleroSlots)) < l {
			his.AcceleroSlots = append(his.AcceleroSlots, make([]bool, l-int64(len(his.AcceleroSlots))+64)...)
		}
		//64 seconds should be enough to predict heart rate
		if l := his.EndTime - his.StartTime + 1; int64(len(his.HeartbeatPrediction)) < l {
			his.HeartbeatPrediction = append(his.HeartbeatPrediction, make([]uint16, l-int64(len(his.HeartbeatPrediction))+64)...)
		}

		// write data to his
		var cursor int64 = in.Time - his.StartTime
		if in.HeartRate != 0 {
			his.HeartBeat[cursor] = in.HeartRate
		}
		if len(in.Accleration1s) > 3*40 {
			his.AcceleroSlots[cursor] = true
			HSet(svcCtx.Ctx, config.Cfg.Rds, "Accelero:"+JwtID, strconv.FormatInt(in.Time, 10), in.Accleration1s)
		}

		HSet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his)
		//predict heart rate if accelerometer data is available
		if len(in.Accleration1s) > 0 {
			go PredicHeartRate(svcCtx.Ctx, cursor, JwtID, in)
		}
		return his, nil
	})

}

func absInt64(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}
