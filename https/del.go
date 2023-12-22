package https

import (
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
)

func (svcCtx *HttpContext) DelHandler() (result interface{}, err error) {
	var (
		jwts      map[string]interface{} = map[string]interface{}{}
		operation string
		rds       *redis.Client
	)
	if rds, err = config.RdsClientByName(svcCtx.RedisDBName); err != nil {
		return nil, err
	}
	svcCtx.MergeJwtField(jwts)

	if operation, err = svcCtx.KeyFieldAtJwt(); err != nil {
		return "", err
	}
	if !permission.IsPermitted(permission.Del, svcCtx.Key, operation) {
		// check operation permission
		return nil, fmt.Errorf(" operation %v not permitted", operation)
	}

	switch svcCtx.Cmd {
	case "HDEL":
		//error if empty Key or Field
		if svcCtx.Field == "" {
			return "false", ErrEmptyKeyOrField
		}
		cmd := rds.HDel(svcCtx.Ctx, svcCtx.Key, svcCtx.Field)
		if err = cmd.Err(); err == nil {
			return "true", nil
		}
		return "false", err
	case "DEL":
		cmd := rds.HDel(svcCtx.Ctx, svcCtx.Key, "del")
		if err = cmd.Err(); err == nil {
			return "true", nil
		}
		return "false", err
	case "ZREM":
		var MemberStr = strings.Split(svcCtx.Req.FormValue("Member"), ",")
		//convert Member to []interface{}
		var Member = make([]interface{}, len(MemberStr))
		for i, v := range MemberStr {
			Member[i] = v
		}
		if err = rds.ZRem(svcCtx.Ctx, svcCtx.Key, Member...).Err(); err == nil {
			return "true", nil
		}
		return "false", err
	case "ZREMRANGEBYSCORE":
		var Min = svcCtx.Req.FormValue("Min")
		var Max = svcCtx.Req.FormValue("Max")
		if err = rds.ZRemRangeByScore(svcCtx.Ctx, svcCtx.Key, Min, Max).Err(); err == nil {
			return "true", nil
		}
		return "false", err
	default:
		return nil, ErrBadCommand
	}

}
