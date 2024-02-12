package api

// Options is parameter to create an API, RPC, or CallAt
type Options struct {
	ServiceName string
	DBName      string
}

// Key purpose of ApiNamed is to allow different API to have the same input type
func OpName(name string) (o Options) {
	o.ServiceName = name
	return
}
func OpDB(DBName string) (o Options) {
	o.DBName = DBName
	return
}
func optionsMerge(options ...Options) (o Options) {
	for _, value := range options {
		if value.ServiceName != "" {
			o.ServiceName = value.ServiceName
		}
		if value.DBName != "" {
			o.DBName = value.DBName
		}
	}
	return o
}
