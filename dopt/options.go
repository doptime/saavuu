// package do stands for data options
package dopt

// DataOptions is parameter to create an API, RPC, or CallAt
type DataOptions struct {
	Key        string
	DataSource string
}

// set a option property
type Setter func(*DataOptions)

// Key purpose of ApiNamed is to allow different API to have the same input type
func Key(key string) Setter {
	return func(opts *DataOptions) {
		opts.Key = key
	}
}
func DataSource(DataSource string) Setter {
	return func(opts *DataOptions) {
		opts.DataSource = DataSource
	}
}
func MergeOptions(options ...Setter) (o *DataOptions) {
	o = &DataOptions{}
	for _, option := range options {
		option(o)
	}
	return o
}
