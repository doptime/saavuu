package https

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/permission"
)

func (svcCtx *HttpContext) GetHandler() (ret interface{}, err error) {
	var (
		jwts     map[string]interface{} = map[string]interface{}{}
		maps     map[string]interface{} = map[string]interface{}{}
		Min, Max string
	)

	svcCtx.MergeJwtField(jwts)

	if len(svcCtx.Key) == 0 {
		return nil, errors.New("no key")
	}
	if act := strings.ToLower(svcCtx.Cmd); !permission.IsGetPermitted(svcCtx.Key, act) {
		// check operation permission
		return nil, fmt.Errorf(" operation %v not permitted", act)
	}

	db := data.Ctx{Ctx: svcCtx.Ctx, Rds: config.DataRds}
	//case Is a member of a set
	switch svcCtx.Cmd {
	case "HGET":
		return ret, db.HGet(svcCtx.Key, svcCtx.Field, &ret)
	case "HGETALL":
		return maps, db.HGetAll(svcCtx.Key, maps)
	case "HMGET":
		return maps, db.HMGET(svcCtx.Key, maps, strings.Split(svcCtx.Field, ",")...)
	case "HKEYS":
		if keys, err := db.HKeys(svcCtx.Key); err != nil {
			return "", err
		} else {
			return json.Marshal(keys)
		}
	case "HEXISTS":
		return db.HExists(svcCtx.Key, svcCtx.Field)
	case "HLEN":
		return db.HLen(svcCtx.Key)
	case "HVALS":
		return db.HVals(svcCtx.Key)
	case "SISMEMBER":
		return db.SIsMember(svcCtx.Key, svcCtx.Req.FormValue("Member"))
	case "TIME":
		pc := data.Ctx{Ctx: svcCtx.Ctx, Rds: config.ParamRds}
		if tm, err := pc.Time(); err != nil {
			return "", err
		} else {
			return tm.UnixMilli(), nil
		}
	case "ZRANGE":
		var (
			start, stop int64         = 0, -1
			result      []interface{} = []interface{}{}
		)
		if start, err = strconv.ParseInt(svcCtx.Req.FormValue("Start"), 10, 64); err != nil {
			return "", errors.New("parse start error:" + err.Error())
		}
		if stop, err = strconv.ParseInt(svcCtx.Req.FormValue("Stop"), 10, 64); err != nil {
			return "", errors.New("parse stop error:" + err.Error())
		}
		// ZRANGE key start stop [WITHSCORES==true]
		if svcCtx.Req.FormValue("WITHSCORES") == "true" {
			var scores []float64
			if scores, err = db.ZRangeWithScores(svcCtx.Key, start, stop, &result); err != nil {
				return "", err
			}
			return json.Marshal(map[string]interface{}{"members": result, "scores": scores})
		}
		// ZRANGE key start stop [WITHSCORES==false]
		if err = db.ZRange(svcCtx.Key, start, stop, &result); err != nil {
			return "", err
		}
		return json.Marshal(result)
	case "ZRANGEBYSCORE":
		var (
			offset, count int64
			scores        []float64
			result        []interface{} = []interface{}{}
		)
		if Min, Max = svcCtx.Req.FormValue("Min"), svcCtx.Req.FormValue("Max"); Min == "" || Max == "" {
			return "", errors.New("no Min or Max")
		}
		//ZRANGEBYSCORE key min max [WITHSCORES==true]
		if svcCtx.Req.FormValue("WITHSCORES") == "true" {
			if scores, err = db.ZRangeByScoreWithScores(svcCtx.Key, &redis.ZRangeBy{Min: Min, Max: Max, Offset: offset, Count: count}, &result); err != nil {
				return "", err
			}
			//marshal result to json
			return json.Marshal(map[string]interface{}{"members": result, "scores": scores})
		}
		//ZRANGEBYSCORE key min max [WITHSCORES==false]
		if err = db.ZRangeByScore(svcCtx.Key, &redis.ZRangeBy{Min: Min, Max: Max, Offset: offset, Count: count}, &result); err != nil {
			return "", err
		}
		return json.Marshal(result)
	case "ZRANK":
		return db.ZRank(svcCtx.Key, svcCtx.Req.FormValue("Member"))
	case "ZCOUNT":
		return db.ZCount(svcCtx.Key, svcCtx.Req.FormValue("Min"), svcCtx.Req.FormValue("Max"))
	case "ZSCORE":
		return db.ZScore(svcCtx.Key, svcCtx.Req.FormValue("Member"))
	}
	return nil, errors.New("unsupported command")

}
