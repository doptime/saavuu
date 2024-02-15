package api

// Options is parameter to create an API, RPC, or CallAt
type Options struct {
	ServiceName string
	DbName      string
}

// set a option property
type Option func(*Options)

// Key purpose of ApiNamed is to allow different API to have the same input type
func OpName(name string) Option {
	return func(opts *Options) {
		opts.ServiceName = name
	}
}
func OpDb(DBName string) Option {
	return func(opts *Options) {
		opts.DbName = DBName
	}
}
func optionsMerge(options ...Option) (o *Options) {
	o = &Options{}
	for _, option := range options {
		option(o)
	}
	return o
}
