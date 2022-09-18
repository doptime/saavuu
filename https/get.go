package https

import (
	"errors"
	"fmt"
	"strings"

	. "saavuu/config"

	"github.com/vmihailenco/msgpack/v5"
)

func replaceUID(scvCtx *HttpContext, s string) (r string, err error) {
	if !strings.Contains(s, "@me") {
		return s, nil
	}
	id, ok := scvCtx.JwtField("id").(string)
	if !ok {
		return s, fmt.Errorf("JWT field id not found")
	}
	s = strings.Replace(s, "@me", id, -1)
	return s, nil
}
func (svcCtx *HttpContext) GetHandler() (ret interface{}, err error) {
	var (
		data         []byte
		resultBytes  []byte                 = []byte{}
		resultString string                 = ""
		result       map[string]interface{} = map[string]interface{}{}
	)
	if svcCtx.Key, err = replaceUID(svcCtx, svcCtx.Key); err != nil {
		return nil, err
	} else if len(svcCtx.Key) == 0 {
		return nil, errors.New("no key")
	}

	if svcCtx.Field, err = replaceUID(svcCtx, svcCtx.Field); err != nil {
		return nil, err
	}
	//return list of keys
	if svcCtx.Field == "" {
		cmd := Cfg.DataRedis.HKeys(svcCtx.Ctx, svcCtx.Key)
		if err = cmd.Err(); err != nil {
			return nil, err
		}
		return msgpack.Marshal(cmd.Val())
	}

	//return item
	cmd := Cfg.DataRedis.HGet(svcCtx.Ctx, svcCtx.Key, svcCtx.Field)
	if data, err = cmd.Bytes(); err != nil {
		return nil, err
	}
	//fill content type, to support binary or json response
	if svcCtx.ResponseContentType != "application/json" {
		if msgpack.Unmarshal(data, &resultBytes) == nil {
			return resultBytes, err
		}
		if msgpack.Unmarshal(data, &resultString) == nil {
			return resultString, err
		}
	}
	if err = msgpack.Unmarshal(data, &result); err != nil {
		return nil, errors.New("unsupported data type")
	}
	//remove fields that not in svc.QueryFields only
	if svcCtx.QueryFields != "" {
		for k := range result {
			if !strings.Contains(svcCtx.QueryFields, k) {
				delete(result, k)
			}
		}
	}
	return result, nil
}
