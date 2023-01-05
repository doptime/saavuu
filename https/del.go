package https

import (
	"errors"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
)

func (svcCtx *HttpContext) DelHandler() (result interface{}, err error) {
	var (
		jwts map[string]interface{} = map[string]interface{}{}
	)
	svcCtx.MergeJwtField(jwts)

	switch svcCtx.Cmd {
	case "HDEL":
		//error if empty Key or Field
		if svcCtx.Key == "" || svcCtx.Field == "" {
			return "false", ErrEmptyKeyOrField
		}
		if !permission.IsPermittedDelOperation(svcCtx.Key, "hdel") {
			return "false", errors.New("permission denied")
		}
		cmd := config.DataRds.HDel(svcCtx.Ctx, svcCtx.Key, svcCtx.Field)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	case "DEL":
		//error if empty Key or Field
		if svcCtx.Key == "" {
			return "false", ErrEmptyKeyOrField
		}
		cmd := config.ParamRds.HDel(svcCtx.Ctx, svcCtx.Key, "del")
		if err = cmd.Err(); err != nil {
			return nil, err
		}
	}

	return "true", nil
}
