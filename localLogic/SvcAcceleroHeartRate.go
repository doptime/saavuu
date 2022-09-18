package localLogic

import (
	"context"
	"fmt"
	"saavuu/config"
	"saavuu/https"
	. "saavuu/redisService"
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
			Start  int64
			End    int64
			HrCnt  int32
			AccCnt int32
		}
	)
	checkInput := func(v interface{}) error {
		in := v.(*Input)
		//validate input
		if in.Accleration1s == nil && in.HeartRate == 0 {
			return fmt.Errorf("ErrInvalidInput")
		}
		if in.Time <= 0 || absInt64(in.Time-time.Now().Unix()) > 3600*24 {
			return fmt.Errorf("ErrInvalidTimeInput")
		}
		if in.Accleration1s != nil && len(in.Accleration1s)%3 != 0 {
			return fmt.Errorf("ErrInvalidAccleration1s")
		}
		return nil
	}

	https.NewLocalService("svc:HeartbeatAcceleros1s", func(svcCtx *https.HttpContext) (data interface{}, err error) {
		//get jwtid
		JwtID, ok := svcCtx.JwtField("id").(string)
		if !ok || len(JwtID) == 0 {
			return nil, https.ErrJWT
		}

		var (
			in            = &Input{}
			rdsKey string = "TrajAccHr:" + JwtID
		)
		if err = https.ToValidStruct(svcCtx.Req, in, checkInput); err != nil {
			return nil, err
		}

		//get previous data
		his := &AcceleroHeartBeat{Start: in.Time, End: in.Time}
		if err = HGet(svcCtx.Ctx, config.Cfg.Rds, rdsKey, "_", &his); err != nil && err != redis.Nil {
			return nil, err
		}
		// save discontinual data to history
		dataDiscontinual := in.Time-his.End > 5*60
		dataTooLong := (his.End - his.Start) > 8*3600
		if dataDiscontinual || dataTooLong {
			//save data to new key,if data less than 1 minute long,discard it
			HSet(svcCtx.Ctx, config.Cfg.Rds, rdsKey, strconv.FormatInt(his.Start, 10), his)
			his.Start = in.Time
			his.End = in.Time
		}

		//make sure capacity of his.HeartBeat is enough, append 64 0s if not
		if in.Time > his.End {
			his.End = in.Time
		}
		// write data to his
		var cursor int64 = in.Time - his.Start
		if in.HeartRate != 0 {
			his.HrCnt++
			go SaveHeartRate(svcCtx.Ctx, JwtID, his.Start, cursor, in.HeartRate)
		}
		if len(in.Accleration1s) > 3*40 {
			his.AccCnt++
			go SaveAccelero(svcCtx.Ctx, JwtID, his.Start, cursor, in.Accleration1s)
		}

		HSet(svcCtx.Ctx, config.Cfg.Rds, rdsKey, "_", &his)
		//predict heart rate if accelerometer data is available
		if len(in.Accleration1s) > 0 {
			go PredicHeartRate(svcCtx.Ctx, JwtID, his.Start, cursor, in.Accleration1s)
		}
		return his, nil
	})
}
func SaveAccelero(ctx context.Context, JwtID string, StartTime int64, Cursor int64, Acceleration1s []int16) (err error) {
	if len(Acceleration1s) == 0 {
		return nil
	}
	var (
		rdsKey string = "TrajAcc:" + JwtID + ":" + strconv.FormatInt(StartTime, 10)
		L      int64
	)
	for L = config.Cfg.Rds.LLen(ctx, rdsKey).Val(); L < Cursor; L++ {
		pip := config.Cfg.Rds.Pipeline()
		pip.RPush(ctx, rdsKey, nil)
		pip.Expire(ctx, rdsKey, -1)
		pip.Exec(ctx)
	}
	if L == Cursor {
		RPush(ctx, config.Cfg.Rds, rdsKey, Acceleration1s)
	} else if Cursor < L {
		LSet(ctx, config.Cfg.Rds, rdsKey, Cursor, Acceleration1s)
	}
	return nil
}
func SaveHeartRate(ctx context.Context, JwtID string, StartTime int64, Cursor int64, HeartRate uint8) (err error) {
	var (
		rdsKey           string = "TrajHr:" + JwtID
		rdsField         string = strconv.FormatInt(StartTime, 10)
		HeartRateHistory []uint8
		Len              int
	)
	if err = HGet(ctx, config.Cfg.Rds, rdsKey, rdsField, &HeartRateHistory); err != nil && err != redis.Nil {
		return err
	}
	//64 seconds should be enough to predict heart rate

	if Len = len(HeartRateHistory); Len == int(Cursor) {
		HeartRateHistory = append(HeartRateHistory, HeartRate)
	} else if Len >= int(Cursor+1) {
		HeartRateHistory[Cursor] = HeartRate
	} else if Len < int(Cursor+1) {
		HeartRateHistory = append(HeartRateHistory, make([]uint8, int(Cursor)+1-Len)...)
		HeartRateHistory[Cursor] = HeartRate
	}
	HSet(ctx, config.Cfg.Rds, rdsKey, rdsField, &HeartRateHistory)
	return nil
}

func PredicHeartRate(ctx context.Context, JwtID string, StartTime int64, Cursor int64, Acceleration1s []int16) (err error) {
	type (
		//heart rate variation predicted by machine learning
		//value /= 256.0 gets the real value of heart rate variation
		svcOut struct {
			Heartbeat float32
		}
	)
	var (
		rdsKey   string = "TrajHrPredict:" + JwtID
		rdsField string = strconv.FormatInt(StartTime, 10)
		//拆分预测数据的理由是，避免写与写心率、加速度数据冲突，因为数据大了以后，io时间会很长，这时候如果存在锁需求，那么冲突几率很高
		HeartBeatPredicted []uint16
		out                *svcOut = &svcOut{}
		Len                int
	)

	paramIn := map[string]interface{}{"UID": JwtID, "Time": Cursor, "Accelero1s": Acceleration1s}
	if err = Call(ctx, config.Cfg.Rds, "svc:PredictHeartBeat", paramIn, out); err != nil {
		return err
	}
	if err = HGet(ctx, config.Cfg.Rds, rdsKey, rdsField, &HeartBeatPredicted); err != nil && err != redis.Nil {
		return err
	}
	heartBeatValue := uint16(out.Heartbeat * 256.0)
	//append directly
	if Len = int(len(HeartBeatPredicted)); Len == int(Cursor) {
		HeartBeatPredicted = append(HeartBeatPredicted, heartBeatValue)
	} else if Len >= int(Cursor+1) {
		HeartBeatPredicted[Cursor] = heartBeatValue
	} else if Len < int(Cursor+1) {
		HeartBeatPredicted = append(HeartBeatPredicted, make([]uint16, int(Cursor)+1-Len)...)
		//write predicted HeartBeat to his
		HeartBeatPredicted[Cursor] = heartBeatValue
	}

	HSet(ctx, config.Cfg.Rds, rdsKey, rdsField, &HeartBeatPredicted)
	fmt.Println("PredictHeartBeat", out.Heartbeat)
	return nil
}

func absInt64(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}
