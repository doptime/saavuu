package https

import (
	"errors"

	"github.com/yangkequn/saavuu/config"
)

func (svcCtx *HttpContext) DelHandler() (result interface{}, err error) {
	var (
		keyWithMyID string
		jwts        map[string]interface{} = map[string]interface{}{}
	)
	svcCtx.MergeJwtField(jwts)

	if keyWithMyID, err = replaceAtUseJwt(svcCtx, jwts, svcCtx.Key); err != nil {
		return nil, err
	} else if keyWithMyID == "" {
		return nil, errors.New("no key")

		//key must contain @me
	} else if keyWithMyID == svcCtx.Key || !(keyWithMyID[0] >= 'A' && keyWithMyID[0] <= 'Z') {
		return nil, errors.New("unauthorized deletion")
	}
	if svcCtx.Field, err = replaceAtUseJwt(svcCtx, jwts, svcCtx.Field); err != nil {
		return nil, err
	}
	if svcCtx.Field == "" {
		cmd := config.ParamRds.Del(svcCtx.Ctx, keyWithMyID)
		if err = cmd.Err(); err != nil {
			return nil, err
		}
		return "{deleted:true,key:" + svcCtx.Key + "} ", nil
	}
	cmd := config.ParamRds.HDel(svcCtx.Ctx, keyWithMyID, svcCtx.Field)
	if err = cmd.Err(); err != nil {
		return nil, err
	}
	return "{deleted:true,key:" + svcCtx.Key + ",field:" + svcCtx.Field + "} ", nil
}
