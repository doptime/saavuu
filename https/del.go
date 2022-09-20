package https

import (
	"errors"

	"github.com/yangkequn/saavuu/config"
)

func (scvCtx *HttpContext) DelHandler() (result interface{}, err error) {
	var (
		keyWithMyID string
	)
	if keyWithMyID, err = replaceUID(scvCtx, scvCtx.Key); err != nil {
		return nil, err
	} else if keyWithMyID == "" {
		return nil, errors.New("no key")

		//key must contain @me
	} else if keyWithMyID == scvCtx.Key || !(keyWithMyID[0] >= 'A' && keyWithMyID[0] <= 'Z') {
		return nil, errors.New("Unauthorized deletion!")
	}
	if scvCtx.Field, err = replaceUID(scvCtx, scvCtx.Field); err != nil {
		return nil, err
	}
	if scvCtx.Field == "" {
		cmd := config.ParamRds.Del(scvCtx.Ctx, keyWithMyID)
		if err = cmd.Err(); err != nil {
			return nil, err
		}
		return "{deleted:true,key:" + scvCtx.Key + "} ", nil
	}
	cmd := config.ParamRds.HDel(scvCtx.Ctx, keyWithMyID, scvCtx.Field)
	if err = cmd.Err(); err != nil {
		return nil, err
	}
	return "{deleted:true,key:" + scvCtx.Key + ",field:" + scvCtx.Field + "} ", nil
}
