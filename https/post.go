package https

import (
	"encoding/json"
	"errors"
	"fmt"
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
		operation string
	)

	if operation, err = svcCtx.KeyFieldAtJwt(); err != nil {
		return "", err
	}
	if !permission.IsPermitted(svcCtx.Key, operation) {
		return "false", ErrOperationNotPermited
	}

	//db := &data.Ctx{Ctx: svcCtx.Ctx, Rds: config.Rds, Key: svcCtx.Key}
	db := data.New[interface{}, interface{}](data.Option.Key(svcCtx.Key))

	//service name is stored in svcCtx.Key
	switch svcCtx.Cmd {
	// all data that appears in the form or body is json format, will be stored in paramIn["JsonPack"]
	// this is used to support 3rd party api
	case "API":
		var (
			paramIn     map[string]interface{} = map[string]interface{}{}
			ServiceName string                 = svcCtx.Key
		)
		svcCtx.MergeJwtField(paramIn)
		//convert query fields to JsonPack. but ignore K field(api name )
		if svcCtx.Req.ParseForm(); len(svcCtx.Req.Form) > 0 {
			for key, value := range svcCtx.Req.Form {
				if paramIn[key] = value[0]; len(value) > 1 {
					paramIn[key] = value // Assign the single value directly
				}
			}
		}
		if msgPack := svcCtx.MsgpackBodyBytes(); len(msgPack) > 0 {
			if err = msgpack.Unmarshal(msgPack, &paramIn); err != nil {
				return nil, fmt.Errorf("msgpack.Unmarshal msgPack error %s", err)
			}
			//paramIn["MsgpackBody"] = msgPack
		} else if jsonBody := svcCtx.JsonBodyBytes(); len(jsonBody) > 0 {
			//convert to msgpack, so that fields can be renamed in ProcessOneJob
			if err = json.Unmarshal(jsonBody, &paramIn); err != nil {
				return nil, fmt.Errorf("msgpack.Unmarshal JsonBody error %s", err)
			}
		}
		return api.CallByHTTP(ServiceName, paramIn, svcCtx.Req)
	case "ZADD":
		var Score float64
		var obj interface{}
		if Score, err = strconv.ParseFloat(svcCtx.Req.FormValue("Score"), 64); err != nil {
			return "false", errors.New("parameter Score shoule be float")
		}
		//unmarshal msgpack
		if MsgPack := svcCtx.MsgpackBodyBytes(); len(MsgPack) == 0 {
			return "false", errors.New("missing MsgPack content")
		} else if err = msgpack.Unmarshal(MsgPack, &obj); err != nil {
			return "false", err
		}
		if err = db.ZAdd(redis.Z{Score: Score, Member: obj}); err != nil {
			return "false", err
		}
		return "true", nil
	default:
		err = ErrBadCommand
	}

	return ret, err
}
