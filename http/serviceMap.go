package http

var ServiceMap map[string]func(svcCtx *HttpContext) (data interface{}, err error) = map[string]func(svcCtx *HttpContext) (data interface{}, err error){}

func NewService(name string, f func(svcCtx *HttpContext) (data interface{}, err error)) {
	ServiceMap[name] = f
}