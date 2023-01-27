package https

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v9"
	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
	"github.com/yangkequn/saavuu/rCtx"
)

func (svcCtx *HttpContext) GetHandler() (ret interface{}, err error) {
	var (
		jwts       map[string]interface{} = map[string]interface{}{}
		maps       map[string]interface{} = map[string]interface{}{}
		_interface interface{}
	)

	svcCtx.MergeJwtField(jwts)

	if len(svcCtx.Key) == 0 {
		return nil, errors.New("no key")
	} else if act := strings.ToLower(svcCtx.Cmd); !permission.IsGetPermitted(svcCtx.Key, act) {
		// check operation permission
		return nil, fmt.Errorf(" operation %v not permitted", act)
	}

	dc := rCtx.DataCtx{Ctx: svcCtx.Ctx, Rds: DataRds}
	//case Is a member of a set
	switch svcCtx.Cmd {
	case "HGET":
		return _interface, dc.HGet(svcCtx.Key, svcCtx.Field, &_interface)
	case "HGETALL":
		return maps, dc.HGetAll(svcCtx.Key, maps)
	case "HMGET":
		return maps, dc.HMGET(svcCtx.Key, maps, strings.Split(svcCtx.Field, ",")...)
	case "HKEYS":
		if keys, err := dc.HKeys(svcCtx.Key); err != nil {
			return "", err
		} else {
			return json.Marshal(keys)
		}
	case "HEXISTS":
		return dc.HExists(svcCtx.Key, svcCtx.Field)
	case "HLEN":
		return dc.HLen(svcCtx.Key)
	case "HVALS":
		return dc.HVals(svcCtx.Key)
	case "SISMEMBER":
		Member := svcCtx.Req.FormValue("Member")
		if Member == "" {
			return "", errors.New("no Member")
		}
		return dc.SIsMember(svcCtx.Key, svcCtx.Field)
	case "ZRANGE":
		var (
			Start, Stop, WITHSCORES string
			start, stop             int64
		)
		if Start = svcCtx.Req.FormValue("Start"); Start == "" {
			return "", errors.New("no Start")
		}
		if Stop = svcCtx.Req.FormValue("Stop"); Stop == "" {
			return "", errors.New("no Stop")
		}
		if WITHSCORES = svcCtx.Req.FormValue("WITHSCORES"); WITHSCORES == "" {
			return "", errors.New("no WITHSCORES")
		}
		if start, err = strconv.ParseInt(Start, 10, 64); err != nil {
			return "", err
		}
		if stop, err = strconv.ParseInt(Stop, 10, 64); err != nil {
			return "", err
		}
		result := []interface{}{}
		if WITHSCORES == "true" {
			cmd := DataRds.ZRangeWithScores(svcCtx.Ctx, svcCtx.Key, start, stop)
			if err = cmd.Err(); err != nil {
				return "", err
			}
			for _, v := range cmd.Val() {
				result = append(result, v.Member)
				result = append(result, v.Score)
			}
		} else {
			cmd := DataRds.ZRange(svcCtx.Ctx, svcCtx.Key, start, stop)
			if err = cmd.Err(); err != nil {
				return "", err
			}
			for _, v := range cmd.Val() {
				result = append(result, v)
			}
		}
		//marshal result to json
		return json.Marshal(result)
	case "ZRANGEBYSCORE":
		var (
			Min, Max, WITHSCORES string
			min, max             float64
		)
		if Min = svcCtx.Req.FormValue("Min"); Min == "" {
			return "", errors.New("no Min")
		}
		if Max = svcCtx.Req.FormValue("Max"); Max == "" {
			return "", errors.New("no Max")
		}
		if WITHSCORES = svcCtx.Req.FormValue("WITHSCORES"); WITHSCORES == "" {
			return "", errors.New("no WITHSCORES")
		}
		result := []interface{}{}
		if WITHSCORES == "true" {
			cmd := DataRds.ZRangeByScoreWithScores(svcCtx.Ctx, svcCtx.Key, &redis.ZRangeBy{
				Min:    Min,
				Max:    Max,
				Offset: 0,
				Count:  0,
			})
			if err = cmd.Err(); err != nil {
				return "", err
			}
			for _, v := range cmd.Val() {
				result = append(result, v.Member)
				result = append(result, v.Score)
			}
		} else {
			cmd := DataRds.ZRangeByScore(svcCtx.Ctx, svcCtx.Key, &redis.ZRangeBy{
				Min:    strconv.FormatFloat(min, 'f', -1, 64),
				Max:    strconv.FormatFloat(max, 'f', -1, 64),
				Offset: 0,
				Count:  0,
			})
			if err = cmd.Err(); err != nil {
				return "", err
			}
			for _, v := range cmd.Val() {
				result = append(result, v)
			}
		}
		//marshal result to json
		return json.Marshal(result)
	case "ZRANK":
		Member := svcCtx.Req.FormValue("Member")
		if Member == "" {
			return "", errors.New("no Member")
		}
		return dc.ZRank(svcCtx.Key, Member)
	}
	return nil, errors.New("unsupported command")

}
