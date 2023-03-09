package permission

var PermittedPutOp map[string]Permission = make(map[string]Permission)

func IsPutPermitted(dataKey string, operation string) (ok bool) {
	return IsPermitted(PermittedPutOp, &permitKeyPut, dataKey, operation)
}
