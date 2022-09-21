package https

import (
	"errors"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/redisContext"

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
	scvCtx.MergeJwtField(paramIn)

	rc := redisContext.RedisContext{Ctx: scvCtx.Ctx, ParamRds: ParamRds, DataRds: DataRds}
	if resultBytes, err = rc.RdsApiBasic(scvCtx.Key, paramIn); err != nil {
		return nil, err
	}

	//fill content type, to support binary or json response
	if err = msgpack.Unmarshal(resultBytes, &result); err != nil {
		return nil, errors.New("unsupported data type")
	} else if err = msgpack.Unmarshal(resultBytes, &responseBytes); err == nil {
		return responseBytes, err
	} else if err = msgpack.Unmarshal(resultBytes, &responseString); err == nil {
		return responseString, err
	}
	return result, nil
}
