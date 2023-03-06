package https

import (
	"errors"

	"github.com/yangkequn/saavuu/api"
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
	if err = api.New(svcCtx.Service).Do(paramIn, &result); err != nil {
		return nil, err
	}

	return result, nil
}
