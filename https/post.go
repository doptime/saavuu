package https

import (
	"errors"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rCtx"

	"github.com/vmihailenco/msgpack/v5"
)

func (svcCtx *HttpContext) PostHandler() (data interface{}, err error) {
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

	pc := rCtx.ParamCtx{Ctx: svcCtx.Ctx, Rds: ParamRds}
	if resultBytes, err = pc.RdsApiBasic(svcCtx.Service, paramIn); err != nil {
		return nil, err
	}

	//if resultBytes is msgpack byte array, return it directly
	if err = msgpack.Unmarshal(resultBytes, &responseBytes); err == nil {
		return responseBytes, err
	}
	//if resultBytes is msgpack string, return it directly
	if err = msgpack.Unmarshal(resultBytes, &responseString); err == nil {
		return responseString, err
	}
	//if resultBytes is msgpack map, return it directly
	if err = msgpack.Unmarshal(resultBytes, &result); err != nil {
		return nil, errors.New("unsupported data type")
	}
	return result, nil
}
