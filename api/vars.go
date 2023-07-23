package api

type InEmpty struct {
}
type InIDPub struct {
	Id  int64 `msgpack:"alias:JWT_id"`
	Pub int64 `msgpack:"alias:JWT_pub"`
}
