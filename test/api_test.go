package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yangkequn/saavuu/api"
)

type InDemo struct {
	Text string
}

var ApiDemo = api.Api(func(InParam *InDemo) (ret string, err error) {
	return "hello world", nil
})

func TestApiDemo(t *testing.T) {

	//create a http context
	var (
		result string
		err    error
	)
	result, err = ApiDemo(&InDemo{Text: "hello"})
	if err != nil {
		t.Error(err)
	} else if result != "hello world" {
		t.Error("result is not hello world")
	}
}

var DemoRpc = api.Rpc[*InDemo, string]()

func TestRPC(t *testing.T) {

	//create a http context
	var result string
	var err error
	result, err = DemoRpc(&InDemo{Text: "hello"})
	if err != nil {
		t.Error(err)
	} else if result != "hello world" {
		t.Error("result is not hello world")
	}

}

func TestCallAt(t *testing.T) {
	var err error

	type InDemo struct {
		Text string
	}
	var demoApi = api.Api(func(InParam *InDemo) (ret string, err error) {
		fmt.Println("hello world called!")
		return "hello world", nil
	})

	var DemoAt = api.CallAt(demoApi)

	//create a http context
	if err = DemoAt(time.Now().Add(time.Second*10), &InDemo{Text: "hello"}); err != nil {
		t.Error(err)
	}
}
