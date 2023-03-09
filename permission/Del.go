package permission

type Permission struct {
	Key       string
	CreateAt  int64
	WhiteList []string
	BlackList []string
}

var PermittedDelOp map[string]Permission = make(map[string]Permission)

func IsDelPermitted(dataKey string, operation string) bool {
	return IsPermitted(PermittedDelOp, &permitKeyDel, dataKey, operation)
}
