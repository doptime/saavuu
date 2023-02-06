package https

import (
	"errors"

	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/config"
)

func (svcCtx *HttpContext) PostHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn map[string]interface{} = map[string]interface{}{}
		result  interface{}            = nil
	)
	if paramIn, err = svcCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	svcCtx.MergeJwtField(paramIn)

	pc := api.Ctx{Ctx: svcCtx.Ctx, Rds: config.ParamRds}
	if err = pc.Do(svcCtx.Service, paramIn, &result); err != nil {
		return nil, err
	}

	return result, nil
}
