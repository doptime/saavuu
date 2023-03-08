package https

import (
	"errors"

	"github.com/yangkequn/saavuu/api"
)

var ErrBadCommand = errors.New("error bad command")

func (svcCtx *HttpContext) PostHandler() (data interface{}, err error) {
	//use remote service map to handle request
	var (
		paramIn map[string]interface{} = map[string]interface{}{}
		result  interface{}            = nil
	)
	if paramIn, err = svcCtx.BodyMessage(); err != nil {
		return nil, errors.New("data error")
	}
	if svcCtx.Cmd == "API" {
		svcCtx.MergeJwtField(paramIn)
		err = api.New(svcCtx.Key).Do(paramIn, &result)
	} else {
		err = ErrBadCommand
	}

	return result, err
}
