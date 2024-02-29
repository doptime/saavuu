package api

// ApiOption is parameter to create an API, RPC, or CallAt
type ApiOption struct {
	Name_       string
	DataSource_ string
}

var Option *ApiOption

// Key purpose of ApiNamed is to allow different API to have the same input type
func (o *ApiOption) Name(apiName string) (out *ApiOption) {
	if out = o; o == Option {
		out = &ApiOption{}
	}
	out.Name_ = apiName
	return out
}

func (o *ApiOption) DataSource(DataSource string) (out *ApiOption) {
	if out = o; o == Option {
		out = &ApiOption{}
	}
	out.DataSource_ = DataSource
	return out
}
