package permission

import cmap "github.com/orcaman/concurrent-map/v2"

var PermittedGetOp cmap.ConcurrentMap[string, Permission] = cmap.New[Permission]()

// Only Batch Get Operation is checked. HGET etc are not checked
func IsGetPermitted(dataKey string, operation string) bool {
	return IsPermitted(PermittedGetOp, &permitKeyGet, dataKey, operation)
}
