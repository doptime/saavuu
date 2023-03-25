package permission

import cmap "github.com/orcaman/concurrent-map/v2"

type Permission struct {
	Key       string
	CreateAt  int64
	WhiteList []string
	BlackList []string
}

var PermittedDelOp cmap.ConcurrentMap[string, Permission] = cmap.New[Permission]()

func IsDelPermitted(dataKey string, operation string) bool {
	return IsPermitted(PermittedDelOp, &permitKeyDel, dataKey, operation)
}
