package https

import (
	"errors"
	"strings"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/redisContext"
	"github.com/yangkequn/saavuu/tools"

	"github.com/vmihailenco/msgpack/v5"
)

func (scvCtx *HttpContext) PutHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn map[string]interface{} = map[string]interface{}{}
		result  map[string]interface{} = map[string]interface{}{}

		resultBytes    []byte = []byte{}
		responseBytes  []byte = []byte{}
		responseString string = ""
	)
	if paramIn, err = scvCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	paramIn["JwtID"] = scvCtx.JwtField("id")
	if id, ok := scvCtx.JwtField("id").(string); ok {
		paramIn["JwtID"] = id
	}

	rc := redisContext.RedisContext{Ctx: scvCtx.Ctx, ParamRds: ParamRds, DataRds: DataRds}
	if resultBytes, err = rc.RdsApiBasic(scvCtx.Key, paramIn); err != nil {
		return nil, err
	}

	//fill content type, to support binary or json response
	if scvCtx.ResponseContentType != "application/json" {
		if err = msgpack.Unmarshal(resultBytes, &responseBytes); err == nil {
			return responseBytes, err
		}
		if err = msgpack.Unmarshal(resultBytes, &responseString); err == nil {
			return responseString, err
		}
	}
	if err = msgpack.Unmarshal(resultBytes, &result); err != nil {
		return nil, errors.New("unsupported data type")
	}
	//remove fields that not in svc.QueryFields only
	if scvCtx.QueryFields != "" {
		for _, k := range tools.MapKeys(result) {
			if !strings.Contains(scvCtx.QueryFields, k) {
				delete(result, k)
			}
		}
	}
	return result, nil
}
