package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/yangkequn/saavuu/api"
)

type Demo1 struct {
	Text string
}

var ApiDemo = api.Api(func(InParam *Demo1) (ret string, err error) {
	return "hello world", nil
})

func TestApiDemo(t *testing.T) {

	//create a http context
	var (
		result string
		err    error
	)
	result, err = ApiDemo(&Demo1{Text: "hello"})
	if err != nil {
		t.Error(err)
	} else if result != "hello world" {
		t.Error("result is not hello world")
	}
}

var DemoRpc = api.Rpc[*Demo1, string]()

func TestRPC(t *testing.T) {

	//create a http context
	var result string
	var err error
	result, err = DemoRpc(&Demo1{Text: "hello"})
	if err != nil {
		time.Sleep(10 * time.Second)
		t.Error(err)
	} else if result != "hello world" {
		t.Error("result is not hello world")
	}

}

func TestCallAt(t *testing.T) {
	var err error
	callAt := api.CallAt(DemoRpc, time.Now().Add(time.Second*10))

	if err = callAt(&Demo1{Text: "TestCallAt 10s later"}); err != nil {
		t.Error(err)
	}
}

func TestHTTPCall(t *testing.T) {
	var (
		body []byte
		err  error
	)
	//create a http context
	rsb, err := http.Get("http://127.0.0.1:8080/api/CallAt?apiName=DemoRpc&timeAt=10s")
	body = make([]byte, 1024)
	rsb.Body.Read(body)

	if err != nil {
		t.Error(err)
	} else if string(body) != "hello world" {
		t.Error("result is not hello world")
	}
}
