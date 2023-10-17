package https

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/permission"
)

var ErrBadCommand = errors.New("error bad command")

func printSertviceCount() {
	fmt.Println("https.ApiServices.Count()", api.ApiServices.Count())
	fmt.Println("https.ApiServices.len", len(api.ApiServices.Items()))
	time.Sleep(time.Second)

}

func (svcCtx *HttpContext) PostHandler() (ret interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn   map[string]interface{} = map[string]interface{}{}
		operation string
		fuc       *api.ApiInfo
		ok        bool
		buf       []byte
	)

	if operation, err = svcCtx.KeyFieldAtJwt(); err != nil {
		return "", err
	}
	if !permission.IsPermitted(permission.Post, svcCtx.Key, operation) {
		return "false", ErrOperationNotPermited
	}

	//db := &data.Ctx{Ctx: svcCtx.Ctx, Rds: config.Rds, Key: svcCtx.Key}
	db := data.New[interface{}, interface{}](svcCtx.Key)

	//service name is stored in svcCtx.Key
	switch svcCtx.Cmd {
	// all data that appears in the form or body is json format, will be stored in paramIn["JsonPack"]
	// this is used to support 3rd party api
	case "JSAPI":
		var (
			fuc *api.ApiInfo
			ok  bool
		)
		//convert query fields to JsonPack. but ignore K field(api name )
		svcCtx.Req.ParseForm()
		if len(svcCtx.Req.Form) > 0 {
			paramIn["JsonPack"], _ = msgpack.Marshal(svcCtx.Req.Form)
		}
		var _api = api.New[interface{}, interface{}](svcCtx.Key)
		//if function is not stored locally, call it remotely (RPC). This is alias microservice mode
		if fuc, ok = api.ApiServices.Get(_api.ServiceName); config.Cfg.RPCFirst || !ok {
			return _api.Do(paramIn)
		}

		//if function is stored locally, call it directly. This is alias monolithic mode
		if buf, err = api.EncodeApiInput(paramIn); err != nil {
			return nil, err
		}
		return fuc.ApiFuncWithMsgpackedParam(buf)
	case "API":
		printSertviceCount()
		if MsgPack, _ := svcCtx.BodyBytes(); len(MsgPack) > 0 {
			paramIn["MsgPack"] = MsgPack
		}
		svcCtx.MergeJwtField(paramIn)
		var _api = api.New[map[string]interface{}, interface{}](svcCtx.Key)
		//if function is not stored locally, call it remotely (RPC). This is alias microservice mode
		if fuc, ok = api.ApiServices.Get(_api.ServiceName); config.Cfg.RPCFirst || !ok {
			return _api.Do(paramIn)
		}
		//if function is stored locally, call it directly. This is alias monolithic mode
		if buf, err = api.EncodeApiInput(paramIn); err != nil {
			return nil, err
		}
		return fuc.ApiFuncWithMsgpackedParam(buf)
	case "ZADD":
		var Score float64
		var obj interface{}
		if Score, err = strconv.ParseFloat(svcCtx.Req.FormValue("Score"), 64); err != nil {
			return "false", errors.New("parameter Score shoule be float")
		}
		//unmarshal msgpack
		if MsgPack, _ := svcCtx.BodyBytes(); len(MsgPack) == 0 {
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
