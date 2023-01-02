package https

import (
	"github.com/yangkequn/saavuu/config"
)

func (svcCtx *HttpContext) DelHandler() (result interface{}, err error) {
	var (
		jwts map[string]interface{} = map[string]interface{}{}
	)
	svcCtx.MergeJwtField(jwts)

	if svcCtx.Field == "" {
		cmd := config.ParamRds.Del(svcCtx.Ctx, svcCtx.Key)
		if err = cmd.Err(); err != nil {
			return nil, err
		}
		return "{deleted:true,key:" + svcCtx.Key + "} ", nil
	}
	cmd := config.ParamRds.HDel(svcCtx.Ctx, svcCtx.Key, svcCtx.Field)
	if err = cmd.Err(); err != nil {
		return nil, err
	}
	return "{deleted:true,key:" + svcCtx.Key + ",field:" + svcCtx.Field + "} ", nil
}
