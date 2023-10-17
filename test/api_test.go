package test

import (
	"testing"

	"github.com/yangkequn/saavuu/api"
)

type InDemo struct {
	Text string
}

var ApiDemo, _ = api.Api(func(InParam *InDemo) (ret string, err error) {
	return "hello world", nil
})

func TestApiDemo(t *testing.T) {
	//create a http context
	var result string
	var err error
	result, err = ApiDemo(&InDemo{Text: "hello"})
	if err != nil {
		t.Error(err)
	}
	if result != "hello world" {
		t.Error("result is not hello world")
	}
}
