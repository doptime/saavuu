package service

import (
	"saavuu/http"
)

var _ bool = http.NewService("SvcAccelero", func(svcCtx *http.HttpContext) (data interface{}, err error) {
	type Input struct {
		Name string
	}
	var i = &Input{}
	if err = http.ToStruct(svcCtx.Req, i); err != nil {
		return nil, err
	}
	return data, nil
})
