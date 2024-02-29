// package do stands for data options
package data

// DataOption is parameter to create an API, RPC, or CallAt
type DataOption struct {
	Key_        string
	DataSource_ string
}

var Option *DataOption

// Key purpose of ApiNamed is to allow different API to have the same input type
func (o *DataOption) Key(key string) (out *DataOption) {
	if out = o; o == Option {
		out = &DataOption{}
	}
	out.Key_ = key
	return out
}
func (o *DataOption) DataSource(dataSource_ string) (out *DataOption) {
	if out = o; o == Option {
		out = &DataOption{}
	}
	out.DataSource_ = dataSource_
	return out
}
