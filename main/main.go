package main

import (
	"github.com/yangkequn/saavuu/https"
)

type TestApi struct {
	ApiBase string
}

func main() {
	https.StartHttp()
}
