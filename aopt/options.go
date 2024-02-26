package aopt

// ApiOptions is parameter to create an API, RPC, or CallAt
type ApiOptions struct {
	Name       string
	DataSource string
}

// set a api option property
// available options are: Name, DS
type Setter func(*ApiOptions)

// Key purpose of ApiNamed is to allow different API to have the same input type
func Name(name string) Setter {
	return func(opts *ApiOptions) {
		opts.Name = name
	}
}
func DataSource(DataSource string) Setter {
	return func(opts *ApiOptions) {
		opts.DataSource = DataSource
	}
}
func MergeOptions(options ...Setter) (o *ApiOptions) {
	o = &ApiOptions{}
	for _, option := range options {
		option(o)
	}
	return o
}
