package api

// Key purpose of ApiNamed is to allow different API to have the same input type
func Name(name string) string {
	return "nm:" + name
}
func DB(name string) string {
	return "db:" + name
}
func optionsDecode(options ...string) (ServiceName, DBName string) {
	for _, value := range options {
		if len(ServiceName) == 0 && len(value) > 3 || value[:3] == "nm:" {
			ServiceName = value[3:]
		}
		if len(DBName) == 0 && len(value) > 3 || value[:3] == "db:" {
			DBName = value[3:]
		}
	}
	return
}
