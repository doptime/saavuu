package https

import (
	"errors"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/redisContext"

	"github.com/vmihailenco/msgpack/v5"
)

func (svcCtx *HttpContext) PutHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn map[string]interface{} = map[string]interface{}{}
		result  map[string]interface{} = map[string]interface{}{}

		resultBytes    []byte = []byte{}
		responseBytes  []byte = []byte{}
		responseString string = ""
	)
	if paramIn, err = svcCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	svcCtx.MergeJwtField(paramIn)

	pc := redisContext.ParamContext{Ctx: svcCtx.Ctx, Rds: ParamRds}
	if resultBytes, err = pc.RdsApiBasic(svcCtx.Service, paramIn); err != nil {
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
