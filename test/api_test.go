package test

import (
	"testing"
	"time"

	"github.com/yangkequn/saavuu/api"
)

type Demo struct {
	Text string
}

var ApiDemo = api.Api(func(InParam *Demo) (ret string, err error) {
	return "hello world", nil
})

func TestApiDemo(t *testing.T) {

	//create a http context
	var (
		result string
		err    error
	)
	result, err = ApiDemo(&Demo{Text: "hello"})
	if err != nil {
		t.Error(err)
	} else if result != "hello world" {
		t.Error("result is not hello world")
	}
}

var DemoRpc = api.Rpc[*Demo, string]()

func TestRPC(t *testing.T) {

	//create a http context
	var result string
	var err error
	result, err = DemoRpc(&Demo{Text: "hello"})
	if err != nil {
		t.Error(err)
	} else if result != "hello world" {
		t.Error("result is not hello world")
	}

}

func TestCallAt(t *testing.T) {
	var err error
	callAt := api.CallAt(DemoRpc, time.Now().Add(time.Second*10))

	if err = callAt(&Demo{Text: "TestCallAt 10s later"}); err != nil {
		t.Error(err)
	}
}
