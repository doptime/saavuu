package https

import (
	"errors"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
)

var ErrEmptyKeyOrField = errors.New("empty key or field")

func (svcCtx *HttpContext) PutHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		result map[string]interface{} = map[string]interface{}{}
		bytes  []byte
	)

	switch svcCtx.Cmd {
	case "HSET":
		//error if empty Key or Field
		if svcCtx.Key == "" || svcCtx.Field == "" {
			return "false", ErrEmptyKeyOrField
		}
		if !permission.IsPutPermitted(svcCtx.Key, "hset") {
			return "false", errors.New("permission denied")
		}
		if bytes, err = svcCtx.MsgpackBody(); err != nil {
			return "false", err
		}
		cmd := DataRds.HSet(svcCtx.Ctx, svcCtx.Key, svcCtx.Field, bytes)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	case "RPUSH":
		//error if empty Key or Field
		if svcCtx.Key == "" {
			return "false", ErrEmptyKeyOrField
		}
		if !permission.IsPutPermitted(svcCtx.Key, "rpush") {
			return "false", errors.New("permission denied")
		}
		if bytes, err = svcCtx.MsgpackBody(); err != nil {
			return "false", err
		}
		cmd := DataRds.RPush(svcCtx.Ctx, svcCtx.Key, bytes)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	}
	return result, nil
}
