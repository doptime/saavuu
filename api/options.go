package api

// ApiOption is parameter to create an API, RPC, or CallAt
type ApiOption struct {
	Name       string
	DataSource string
}

var Option *ApiOption

// Key purpose of ApiNamed is to allow different API to have the same input type
func (o *ApiOption) WithName(apiName string) (out *ApiOption) {
	if out = o; o == Option {
		out = &ApiOption{}
	}
	out.Name = apiName
	return out
}

func (o *ApiOption) WithDataSource(DataSource string) (out *ApiOption) {
	if out = o; o == Option {
		out = &ApiOption{}
	}
	out.DataSource = DataSource
	return out
}
