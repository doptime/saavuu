package http

import (
	"errors"
	. "saavuu/config"
	. "saavuu/redis"
	"saavuu/tools"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

func (scvCtx *HttpContext) PutHandler() (data interface{}, err error) {
	//use local service map to handle request
	if fun, ok := ServiceMap[scvCtx.Service]; ok {
		return fun(scvCtx)
	}
	//use remote service map to handle request
	var (
		paramIn map[string]interface{} = map[string]interface{}{}
		result  map[string]interface{} = map[string]interface{}{}

		resultBytes  []byte = []byte{}
		resultString string = ""
	)
	if paramIn, err = scvCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	if resultBytes, err = DoBasic(scvCtx.Ctx, Cfg.Rds, scvCtx.Key, paramIn); err != nil {
		return nil, err
	}

	//fill content type, to support binary or json response
	if scvCtx.ExpectedReponseType != "application/json" {
		if err = msgpack.Unmarshal(resultBytes, &resultBytes); err == nil {
			return resultBytes, err
		}
		if err = msgpack.Unmarshal(resultBytes, &resultString); err == nil {
			return resultString, err
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
