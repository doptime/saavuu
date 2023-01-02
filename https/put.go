package https

import (
	"errors"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
)

func (svcCtx *HttpContext) PutHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn map[string]interface{} = map[string]interface{}{}
		result  map[string]interface{} = map[string]interface{}{}
		bytes   []byte
	)
	if paramIn, err = svcCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	svcCtx.MergeJwtField(paramIn)

	switch svcCtx.Cmd {
	case "HSET":
		//error if empty Key or Field
		if svcCtx.Key == "" || svcCtx.Field == "" {
			return "false", errors.New("empty key or field")
		}
		if !permission.IsPermittedPutOperation(svcCtx.Key, svcCtx.Field) {
			return "false", errors.New("permission denied")
		}
		if bytes, err = svcCtx.BodyBytes(); err != nil {
			return "false", err
		}
		cmd := DataRds.HSet(svcCtx.Ctx, svcCtx.Key, svcCtx.Field, bytes)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	case "HDEL":
		//error if empty Key or Field
		if svcCtx.Key == "" || svcCtx.Field == "" {
			return "false", errors.New("empty key or field")
		}
		if !permission.IsPermittedPutOperation(svcCtx.Key, svcCtx.Field) {
			return "false", errors.New("permission denied")
		}
		cmd := DataRds.HDel(svcCtx.Ctx, svcCtx.Key, svcCtx.Field)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	}
	return result, nil
}
