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
func (svcCtx *ServiceContext) getHandler() (data []byte, err error) {
	var (
		resultBytes  []byte                 = []byte{}
		resultString string                 = ""
		result       map[string]interface{} = map[string]interface{}{}
	)
	if svcCtx.Key, err = svcCtx.replaceUID(svcCtx.Key); err != nil {
		return nil, err
	} else if len(svcCtx.Key) == 0 {
		return nil, errors.New("no key")
	}

	if svcCtx.Field, err = svcCtx.replaceUID(svcCtx.Field); err != nil {
		return nil, err
	}
	//return list of keys
	if svcCtx.Field == "" {
		cmd := Config.rds.HKeys(svcCtx.ctx, svcCtx.Key)
		if err = cmd.Err(); err != nil {
			return nil, err
		}
		return msgpack.Marshal(cmd.Val())
	}

	//return item
	cmd := Config.rds.HGet(svcCtx.ctx, svcCtx.Key, svcCtx.Field)
	if data, err = cmd.Bytes(); err != nil {
		return nil, err
	}
	//fill content type, to support binary or json response
	if svcCtx.ExpectedReponseType != "application/json" {
		if msgpack.Unmarshal(data, &resultBytes) == nil {
			return resultBytes, err
		}
		if msgpack.Unmarshal(data, &resultString) == nil {
			return []byte(resultString), err
		}
	}
	if err = msgpack.Unmarshal(data, &result); err == nil {
		//remove fields that not in svc.QueryFields only
		if svcCtx.QueryFields != "" {
			for k := range result {
				if !strings.Contains(svcCtx.QueryFields, k) {
					delete(result, k)
				}
			}
		}
		return json.Marshal(result)
	}
	return nil, errors.New("unsupported data type")
}
