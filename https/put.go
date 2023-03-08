package https

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
)

var ErrEmptyKeyOrField = errors.New("empty key or field")
var ErrOperationNotPermited = errors.New("operation permission denied")

func (svcCtx *HttpContext) PutHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		bytes     []byte
		operation string = strings.ToLower(svcCtx.Cmd)
	)

	if strings.Contains(svcCtx.Field, "@") {
		if err := svcCtx.ParseJwtToken(); err != nil {
			return "false", fmt.Errorf("parse JWT token error: %v", err)
		}
		if operation, err = permission.IsPermittedPutField(operation, &svcCtx.Field, svcCtx.jwtToken); err != nil {
			return "false", ErrOperationNotPermited
		}
	}
	if !permission.IsPutPermitted(svcCtx.Key, operation) {
		return "false", ErrOperationNotPermited
	}

	switch svcCtx.Cmd {
	case "HSET":
		//error if empty Key or Field
		if svcCtx.Key == "" || svcCtx.Field == "" {
			return "false", ErrEmptyKeyOrField
		}
		if bytes, err = svcCtx.MsgpackBody(); err != nil {
			return "false", err
		}
		cmd := config.DataRds.HSet(svcCtx.Ctx, svcCtx.Key, svcCtx.Field, bytes)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	case "RPUSH":
		//error if empty Key or Field
		if svcCtx.Key == "" {
			return "false", ErrEmptyKeyOrField
		}
		if bytes, err = svcCtx.MsgpackBody(); err != nil {
			return "false", err
		}
		cmd := config.DataRds.RPush(svcCtx.Ctx, svcCtx.Key, bytes)
		if err = cmd.Err(); err != nil {
			return "false", err
		}
		return "true", nil
	default:
		return nil, ErrBadCommand
	}
}
