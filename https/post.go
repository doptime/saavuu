package https

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/permission"
)

var ErrBadCommand = errors.New("error bad command")

func (svcCtx *HttpContext) PostHandler() (ret interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn   map[string]interface{} = map[string]interface{}{}
		operation string                 = strings.ToLower(svcCtx.Cmd)
	)
	if paramIn, err = svcCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	if strings.Contains(svcCtx.Field, "@") {
		if err := svcCtx.ParseJwtToken(); err != nil {
			return "false", fmt.Errorf("parse JWT token error: %v", err)
		}
		if operation, err = permission.IsPermittedField(operation, &svcCtx.Field, svcCtx.jwtToken); err != nil {
			return "false", ErrOperationNotPermited
		}
	}
	if !permission.IsPutPermitted(svcCtx.Key, operation) {
		return "false", ErrOperationNotPermited
	}

	//db := &data.Ctx{Ctx: svcCtx.Ctx, Rds: config.DataRds, Key: svcCtx.Key}
	db := data.New(svcCtx.Key)

	if svcCtx.Cmd == "API" {
		svcCtx.MergeJwtField(paramIn)
		err = api.New(svcCtx.Key).Do(paramIn, &ret)
	} else if svcCtx.Cmd == "ZADD" {
		var (
			ScoreStr string
			Score    float64
		)
		svcCtx.MergeJwtField(paramIn)
		//for each field, check if it contains @, replace it with jwt field
		if ScoreStr = svcCtx.Req.FormValue("Score"); ScoreStr == "" {
			return "false", errors.New("score is empty")
		}
		if Score, err = strconv.ParseFloat(ScoreStr, 64); err != nil {
			return "false", errors.New("score is not float")
		}
		if err = db.ZAdd(redis.Z{Score: Score, Member: paramIn["MsgPack"]}); err != nil {
			return "false", err
		}
		return "true", nil
	} else {
		err = ErrBadCommand
	}

	return ret, err
}
