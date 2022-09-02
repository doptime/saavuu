package service

import (
	"fmt"
	"saavuu/config"
	"saavuu/http"
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
			HeartBeat     uint8
			//use client's time
			Time int64
		}
		HeartBeatPredicted struct {
			Heartbeat uint8
		}
		AcceleroHeartBeat struct {
			//use client's time
			StartTime int64
			EndTime   int64
			//relative time to startTime
			AcceleroSlots []uint32
			//put heartbeats here to easily return to client
			HeartBeat           []uint8
			HeartbeatPrediction []uint8
		}
	)
	checkInput := func(v interface{}) error {
		in := v.(Input)
		//validate input
		if in.Accleration1s == nil && in.HeartBeat == 0 {
			return fmt.Errorf("ErrInvalidInput")
		}
		if in.Time <= 0 || absInt64(in.Time-time.Now().Unix()) > 86400 {
			return fmt.Errorf("ErrInvalidTimeInput")
		}
		if in.Accleration1s != nil && len(in.Accleration1s)%3 != 0 {
			return fmt.Errorf("ErrInvalidAccleration1s")
		}
		return nil
	}
	http.NewService("svc:HeartbeatAcceleros1s", func(svcCtx *http.HttpContext) (data interface{}, err error) {
		//get jwtid
		JwtID, ok := svcCtx.JwtField("id").(string)
		if !ok {
			return nil, http.ErrJWT
		}

		var (
			in = &Input{}
		)
		if err = http.ToValidStruct(svcCtx.Req, in, checkInput); err != nil {
			return nil, err
		}

		fmt.Println("lengthof data is ", len(in.Accleration1s)/3)

		//get previous data
		his := &AcceleroHeartBeat{StartTime: in.Time, EndTime: in.Time}
		if err = HGet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his); err != nil && err != redis.Nil {
			return nil, err
		}
		// save discontinual data to history
		dataDiscontinual := in.Time-his.EndTime > 60
		dataTooLong := (his.EndTime - his.StartTime) > 24*3600
		if dataDiscontinual || dataTooLong {
			//save data to new key,if data less than 1 minute long,discard it
			HSet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, strconv.FormatInt(his.StartTime, 10), his)
			his.StartTime = in.Time
			his.EndTime = in.Time
			his.AcceleroSlots = make([]uint32, 64)
			his.HeartBeat = make([]uint8, 64)
			his.HeartbeatPrediction = make([]uint8, 64)
		}

		//make sure capacity of his.HeartBeat is enough, append 64 0s if not
		if l := his.EndTime - his.StartTime + 1; int64(len(his.HeartBeat)) < l {
			his.HeartBeat = append(his.HeartBeat, make([]uint8, l-int64(len(his.HeartBeat))+64)...)
		}
		if l := his.EndTime - his.StartTime + 1; int64(len(his.HeartbeatPrediction)) < l {
			his.HeartbeatPrediction = append(his.HeartbeatPrediction, make([]uint8, l-int64(len(his.HeartbeatPrediction))+64)...)
		}
		if l := his.EndTime - his.StartTime + 1; int64(len(his.AcceleroSlots)) < l {
			his.AcceleroSlots = append(his.AcceleroSlots, make([]uint32, l-int64(len(his.AcceleroSlots))+64)...)
		}

		// write data to his
		his.EndTime = in.Time
		var currentSlot uint32 = uint32(his.EndTime - his.StartTime)
		if in.HeartBeat != 0 {
			his.HeartBeat[currentSlot] = in.HeartBeat
		}
		if in.Accleration1s != nil {
			his.AcceleroSlots[currentSlot] = currentSlot
			HSet(svcCtx.Ctx, config.Cfg.Rds, "Accelero:"+JwtID, strconv.Itoa(int(currentSlot)), in.Accleration1s)
		}

		if len(svcCtx.QueryFields) == 0 {
			// convert to HeartRate1s
			out := &HeartBeatPredicted{}
			if err = Do(svcCtx.Ctx, config.Cfg.Rds, "svc:Acceleros1sToHeartBeat", in, out); err != nil {
				HSet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his)
				return nil, err
			}
			//write predicted HeartBeat to his
			his.HeartbeatPrediction[currentSlot] = out.Heartbeat
		}
		HSet(svcCtx.Ctx, config.Cfg.Rds, "AcceleroHeartbeat:"+JwtID, "_", &his)
		return his, nil
	})
}

func absInt64(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}
