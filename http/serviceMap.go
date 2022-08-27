package http

var ServiceMap map[string]func(svcCtx *HttpContext) (data interface{}, err error) = map[string]func(svcCtx *HttpContext) (data interface{}, err error){}
