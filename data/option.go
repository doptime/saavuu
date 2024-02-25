package data

import "context"

// Options is parameter to create an API, RPC, or CallAt
type Options struct {
	Key            string
	DataSourceName string
	Ctx            *context.Context
}

// set a option property
type With func(*Options)

// Key purpose of ApiNamed is to allow different API to have the same input type
func WithKey(key string) With {
	return func(opts *Options) {
		opts.Key = key
	}
}
func WithDS(DataSourceName string) With {
	return func(opts *Options) {
		opts.DataSourceName = DataSourceName
	}
}
func WithContext(ctx *context.Context) With {
	return func(opts *Options) {
		opts.Ctx = ctx
	}
}
func mergeOptions(options ...With) (o *Options) {
	o = &Options{}
	for _, option := range options {
		option(o)
	}
	if o.Ctx == nil {
		*o.Ctx = context.Background()
	}
	return o
}
