package https

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/permission"
)

func mapConvertWithKeyFromInterfaceToString(m map[interface{}]interface{}) (m2 map[string]interface{}, err error) {
	var bytes []byte
	m2 = make(map[string]interface{})
	for k, v := range m {
		//marshal to string,using json
		if reflect.TypeOf(k).Kind() == reflect.String {
			m2[k.(string)] = v
			continue
		} else if bytes, err = json.Marshal(k); err != nil {
			return nil, err
		}
		m2[string(bytes)] = v
	}
	return m2, nil
}

func (svcCtx *HttpContext) GetHandler() (ret interface{}, err error) {
	var (
		jwts                    map[string]interface{} = map[string]interface{}{}
		Min, Max                string
		tm                      time.Time
		map_interface_interface map[interface{}]interface{}
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
		if err := db.HGetAll(svcCtx.Key, &map_interface_interface); err != nil {
			return nil, err
		}
		return mapConvertWithKeyFromInterfaceToString(map_interface_interface)
	case "HMGET":
		if err = db.HMGET(svcCtx.Key, strings.Split(svcCtx.Field, ","), &map_interface_interface); err != nil {
			return nil, err
		}
		return mapConvertWithKeyFromInterfaceToString(map_interface_interface)
	case "HKEYS":
		var keys []string
		if err := db.HKeys(svcCtx.Key, &keys); err != nil {
			return "", err
		} else {
			return json.Marshal(keys)
		}
	case "HEXISTS":
		return db.HExists(svcCtx.Key, svcCtx.Field)
	case "HLEN":
		return db.HLen(svcCtx.Key)
	case "HVALS":
		var values []interface{}
		if err = db.HVals(svcCtx.Key, &values); err != nil {
			return "", err
		}
		return values, nil
	case "SISMEMBER":
		return db.SIsMember(svcCtx.Key, svcCtx.Req.FormValue("Member"))
	case "TIME":
		pc := data.Ctx{Ctx: svcCtx.Ctx, Rds: config.ParamRds}
		if tm, err = pc.Time(); err != nil {
			return "", err
		}
		return tm.UnixMilli(), nil
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
