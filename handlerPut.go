package saavuu

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/vmihailenco/msgpack"
)

func (scvCtx *ServiceContext) putHandler() (data []byte, err error) {
	var (
		dataIn  string
		paramIn map[string]interface{} = map[string]interface{}{}
		ok      bool                   = false

		resultBytes  []byte                 = []byte{}
		resultString string                 = ""
		result       map[string]interface{} = map[string]interface{}{}
		contentType  string                 = "application/json"
	)
	if dataIn, ok = scvCtx.Data(); !ok {
		return nil, errors.New("data error")
	}
	if err = json.Unmarshal([]byte(dataIn), &paramIn); err != nil {
		return nil, err
	}
	if resultBytes, err = RedisDo(scvCtx.ctx, Config.rds, scvCtx.Key, paramIn); err != nil {
		return nil, err
	}

	//fill content type, to support binary or json response
	if _Type := scvCtx.req.Header.Get("Content-Type"); _Type != "application/json" {
		scvCtx.rsb.Header().Set("Content-Type", _Type)
		if err = msgpack.Unmarshal(resultBytes, &resultBytes); err == nil {
			return resultBytes, err
		}
		if err = msgpack.Unmarshal(resultBytes, &resultString); err == nil {
			return []byte(resultString), err
		}
	}
	scvCtx.rsb.Header().Set("Content-Type", contentType)
	if err = msgpack.Unmarshal(data, &result); err == nil {
		//remove fields that not in svc.QueryFields only
		if scvCtx.QueryFields != "" {
			for k := range result {
				if !strings.Contains(scvCtx.QueryFields, k) {
					delete(result, k)
				}
			}
		}
		return json.Marshal(result)
	}
	return nil, errors.New("unsupported data type")
}
