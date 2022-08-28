package service

import (
	"saavuu/http"
)

type Input struct {
	Name string
}

var _ = http.NewService("SvcDemo", func(svcCtx *http.HttpContext) (data interface{}, err error) {
	var i = &Input{}
	if err = http.ToStruct(svcCtx.Req, i); err != nil {
		return nil, err
	}
	// your logic here

	return data, nil
})
