package service

import (
	"saavuu/http"
)

var SvcDemo = func() func(svcCtx *http.HttpContext) (data interface{}, err error) {
	var fn = func(svcCtx *http.HttpContext) (data interface{}, err error) {
		return []byte("Hello " + svcCtx.Key), nil
	}
	http.ServiceMap["SvcDemo"] = fn
	return fn
}
