package https

import (
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/permission"
)

var ErrBadCommand = errors.New("error bad command")

func (svcCtx *HttpContext) PostHandler() (ret interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn   map[string]interface{} = map[string]interface{}{}
		operation string
	)
	if paramIn, err = svcCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	if operation, err = svcCtx.KeyFieldAtJwt(); err != nil {
		return "", err
	}
	if !permission.IsPutPermitted(svcCtx.Key, operation) {
		return "false", ErrOperationNotPermited
	}

	//db := &data.Ctx{Ctx: svcCtx.Ctx, Rds: config.DataRds, Key: svcCtx.Key}
	db := data.New[interface{}](svcCtx.Key)

	if svcCtx.Cmd == "API" {
		svcCtx.MergeJwtField(paramIn)
		err = api.New[map[string]interface{}](svcCtx.Key).Do(paramIn, &ret)
	} else if svcCtx.Cmd == "ZADD" {
		var Score float64
		var obj interface{}
		var ok bool
		if Score, err = strconv.ParseFloat(svcCtx.Req.FormValue("Score"), 64); err != nil {
			return "false", errors.New("parameter Score shoule be float")
		}
		//unmarshal msgpack
		if _, ok = paramIn["MsgPack"]; !ok {
			return "false", errors.New("missing MsgPack content")
		}
		if err = msgpack.Unmarshal(paramIn["MsgPack"].([]byte), &obj); err != nil {
			return "false", err
		}
		if err = db.ZAdd(redis.Z{Score: Score, Member: obj}); err != nil {
			return "false", err
		}
		return "true", nil
	} else {
		err = ErrBadCommand
	}

	return ret, err
}
