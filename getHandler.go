package saavuu

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack"
)

func (scvCtx *ServiceContext) replaceUID(s string) (r string, err error) {
	if !strings.Contains(s, "@me") {
		return s, nil
	}
	id, ok := scvCtx.JwtField("id")
	if !ok {
		return s, fmt.Errorf("no id in jwt")
	}
	s = strings.Replace(s, "@me", id, -1)
	return s, nil
}
func (scvCtx *ServiceContext) getHandler() (data []byte, err error) {
	var (
		resultBytes  []byte                 = []byte{}
		resultString string                 = ""
		result       map[string]interface{} = map[string]interface{}{}
		contentType  string                 = "application/json"
	)
	if scvCtx.Key, err = scvCtx.replaceUID(scvCtx.Key); err != nil {
		return nil, err
	} else if len(scvCtx.Key) == 0 {
		return nil, errors.New("no key")
	}

	if scvCtx.Field, err = scvCtx.replaceUID(scvCtx.Field); err != nil {
		return nil, err
	}
	//return list of keys
	if scvCtx.Field == "" {
		cmd := Config.rds.HKeys(scvCtx.ctx, scvCtx.Key)
		if err = cmd.Err(); err != nil {
			return nil, err
		}
		return msgpack.Marshal(cmd.Val())
	}

	//return item
	cmd := Config.rds.HGet(scvCtx.ctx, scvCtx.Key, scvCtx.Field)
	if data, err = cmd.Bytes(); err != nil {
		return nil, err
	}
	//fill content type, to support binary or json response
	if _Type := scvCtx.req.Header.Get("Content-Type"); _Type != "application/json" {
		scvCtx.rsb.Header().Set("Content-Type", _Type)
		if msgpack.Unmarshal(data, resultBytes) == nil {
			return resultBytes, err
		}
		if msgpack.Unmarshal(data, resultString) == nil {
			return []byte(resultString), err
		}
	}
	scvCtx.rsb.Header().Set("Content-Type", contentType)
	if err = msgpack.Unmarshal(data, result); err == nil {
		//remove fields that not in svc.QueryFields only
		if scvCtx.QueryFields != "" {
			for k, _ := range result {
				if !strings.Contains(scvCtx.QueryFields, k) {
					delete(result, k)
				}
			}
		}
		return json.Marshal(result)
	}
	return nil, errors.New("unsupported data type")
}
