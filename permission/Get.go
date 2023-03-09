package permission

var PermittedGetOp map[string]Permission = make(map[string]Permission)

// Only Batch Get Operation is checked. HGET etc are not checked
func IsGetPermitted(dataKey string, operation string) bool {
	return IsPermitted(PermittedGetOp, &permitKeyGet, dataKey, operation)
}
