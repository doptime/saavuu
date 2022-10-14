package https

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v9"
	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rCtx"

	"github.com/vmihailenco/msgpack/v5"
)

func replaceAtUseJwt(scvCtx *HttpContext, jwts map[string]interface{}, s string) (r string, err error) {
	var (
		jwtValue    interface{}
		ok          bool
		jwtFieldStr string
	)
	for i := strings.Index(s, "@"); i >= 0; i = strings.Index(s, "@") {
		ie := i + 1
		for ; ie < len(s) && ((s[ie] >= '0' && s[ie] <= '9') || (s[ie] >= 'a' && s[ie] <= 'z') || (s[ie] >= 'A' && s[ie] <= 'Z')); ie++ {
		}
		if ie == i+1 {
			return "", errors.New("invalid field")
		}
		label := s[i+1 : ie]
		if jwtValue, ok = jwts["JWT_"+label]; !ok {
			return "", errors.New("invalid jwt field")
		}
		if jwtFieldStr, ok = jwtValue.(string); !ok {
			return "", errors.New("invalid jwt field")
		}
		s = s[:i] + jwtFieldStr + s[ie:]
	}
	return s, nil
}
func (svcCtx *HttpContext) GetHandler() (ret interface{}, err error) {
	var (
		jwts         map[string]interface{} = map[string]interface{}{}
		maps         map[string]interface{} = map[string]interface{}{}
		data         []byte
		resultBytes  []byte = []byte{}
		resultString string = ""
	)

	svcCtx.MergeJwtField(jwts)

	if svcCtx.Key, err = replaceAtUseJwt(svcCtx, jwts, svcCtx.Key); err != nil {
		return nil, err
	} else if len(svcCtx.Key) == 0 {
		return nil, errors.New("no key")
	}

	if svcCtx.Field, err = replaceAtUseJwt(svcCtx, jwts, svcCtx.Field); err != nil {
		return nil, err
	}

	//check auth. only Key start with upper case are allowed to access
	if len(svcCtx.Key) <= 0 || !(svcCtx.Key[0] >= 'A' && svcCtx.Key[0] <= 'Z') {
		return nil, errors.New("private Key")
	}
	//case Is a member of a set
	switch svcCtx.Cmd {
	case "HGET":
		cmd := DataRds.HGet(svcCtx.Ctx, svcCtx.Key, svcCtx.Field)
		if data, err = cmd.Bytes(); err != nil {
			return "", err
		}
		//fill content type, to support binary or json response
		if svcCtx.ResponseContentType != "application/json" {
			if msgpack.Unmarshal(data, &resultBytes) == nil {
				return resultBytes, err
			}
			if msgpack.Unmarshal(data, &resultString) == nil {
				return resultString, err
			}
		}

		var _v interface{}
		if err = msgpack.Unmarshal(data, &_v); err != nil {
			return nil, errors.New("unsupported data type")
		}
		return json.Marshal(_v)
	case "HGETALL":
		cmd := DataRds.HGetAll(svcCtx.Ctx, svcCtx.Key)
		if err = cmd.Err(); err != nil {
			return "", err
		}
		for k, v := range cmd.Val() {
			var _v interface{}
			if err = msgpack.Unmarshal([]byte(v), &_v); err != nil {
				continue
			}
			maps[k] = _v
		}
		return json.Marshal(maps)

	case "HMGET":
		cmd := DataRds.HMGet(svcCtx.Ctx, svcCtx.Key, strings.Split(svcCtx.Field, ",")...)
		if err = cmd.Err(); err != nil {
			return "", err
		}
		for i, v := range cmd.Val() {
			if v == nil {
				continue
			}
			var _v interface{}
			if err = msgpack.Unmarshal([]byte(v.(string)), &_v); err != nil {
				continue
			}
			maps[strings.Split(svcCtx.Field, ",")[i]] = _v
		}
		return json.Marshal(maps)
	case "HKEYS":
		cmd := DataRds.HKeys(svcCtx.Ctx, svcCtx.Key)
		if err = cmd.Err(); err != nil {
			return "", err
		}
		return json.Marshal(cmd.Val())
	case "HEXISTS":
		dc := rCtx.DataCtx{Ctx: svcCtx.Ctx, Rds: DataRds}
		if ok := dc.HExists(svcCtx.Key, svcCtx.Field); ok {
			return "true", nil
		}
		return "false", nil
	case "HLEN":
		cmd := DataRds.HLen(svcCtx.Ctx, svcCtx.Key)
		if err = cmd.Err(); err != nil {
			return "", err
		}
		return strconv.FormatInt(cmd.Val(), 10), nil
	case "HVALS":
		cmd := DataRds.HVals(svcCtx.Ctx, svcCtx.Key)
		if err = cmd.Err(); err != nil {
			return "", err
		}
		result := []interface{}{}
		for _, v := range cmd.Val() {

			var _v interface{}
			if err = msgpack.Unmarshal([]byte(v), &_v); err != nil {
				continue
			}
			result = append(result, _v)
		}
		return json.Marshal(result)
	case "SISMEMBER":
		Member := svcCtx.Req.FormValue("Member")
		if Member == "" {
			return "", errors.New("no Member")
		}
		dc := rCtx.DataCtx{Ctx: svcCtx.Ctx, Rds: DataRds}
		if ok := dc.SIsMember(svcCtx.Key, svcCtx.Field); ok {
			return "true", nil
		}
		return "false", nil
	case "ZRANGE":
		var (
			Start, Stop, WITHSCORES string
			start, stop             int64
			withScores              bool = false
		)
		if Start := svcCtx.Req.FormValue("Start"); Start == "" {
			return "", errors.New("no Start")
		}
		if Stop := svcCtx.Req.FormValue("Stop"); Stop == "" {
			return "", errors.New("no Stop")
		}
		if WITHSCORES := svcCtx.Req.FormValue("WITHSCORES"); WITHSCORES == "" {
			return "", errors.New("no WITHSCORES")
		}
		if start, err = strconv.ParseInt(Start, 10, 64); err != nil {
			return "", err
		}
		if stop, err = strconv.ParseInt(Stop, 10, 64); err != nil {
			return "", err
		}
		if WITHSCORES == "true" {
			withScores = true
		}
		result := []interface{}{}
		if withScores {
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
			withScores           bool = false
		)
		if Min = svcCtx.Req.FormValue("Min"); Min == "" {
			return "", errors.New("no Min")
		}
		if Max = svcCtx.Req.FormValue("Max"); Max == "" {
			return "", errors.New("no Max")
		}
		if WITHSCORES := svcCtx.Req.FormValue("WITHSCORES"); WITHSCORES == "" {
			return "", errors.New("no WITHSCORES")
		}
		if WITHSCORES == "true" {
			withScores = true
		}
		result := []interface{}{}
		if withScores {
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
		cmd := DataRds.ZRank(svcCtx.Ctx, svcCtx.Key, Member)
		if err = cmd.Err(); err != nil {
			return "", err
		}
		return strconv.FormatInt(cmd.Val(), 10), nil

	}
	return nil, errors.New("unsupported command")

}
