package service

import (
	"saavuu/http"
)

type Input struct {
	Name string
}

var SvcDemo = http.NewService("SvcDemo", func(svcCtx *http.HttpContext) (data interface{}, err error) {
	var i = &Input{}
	if err = http.ToStruct(svcCtx.Req, i); err != nil {
		return nil, err
	}
	return data, nil
})
