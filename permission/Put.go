package permission

import cmap "github.com/orcaman/concurrent-map/v2"

var PermittedPutOp cmap.ConcurrentMap[string, Permission] = cmap.New[Permission]()

func IsPutPermitted(dataKey string, operation string) (ok bool) {
	return IsPermitted(PermittedPutOp, &permitKeyPut, dataKey, operation)
}
