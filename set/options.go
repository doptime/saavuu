package set

// ApiSetting is parameter to create an API, RPC, or CallAt
type ApiSetting struct {
	Name       string
	DataSource string
}

// set a api option property
// available options are: Name, DS
type Api func(*ApiSetting)

// Key purpose of ApiNamed is to allow different API to have the same input type
func Name(name string) Api {
	return func(opts *ApiSetting) {
		opts.Name = name
	}
}
func DataSource(DataSourceName string) Api {
	return func(opts *ApiSetting) {
		opts.DataSource = DataSourceName
	}
}
func Merge(options ...Api) (o *ApiSetting) {
	o = &ApiSetting{}
	for _, option := range options {
		option(o)
	}
	return o
}
