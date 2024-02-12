package test

import (
	"testing"

	"github.com/yangkequn/saavuu/api"
)

func TestApiDemo(t *testing.T) {

	type InDemo struct {
		Text string
	}

	var ApiDemo = api.Api(func(InParam *InDemo) (ret string, err error) {
		return "hello world", nil
	})

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
