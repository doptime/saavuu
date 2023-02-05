package https

import (
	"errors"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rCtx"
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

	pc := rCtx.ApiCtx{Ctx: svcCtx.Ctx, Rds: ParamRds}
	if err = pc.Api(svcCtx.Service, paramIn, &result, 0); err != nil {
		return nil, err
	}

	return result, nil
}
