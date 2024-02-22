package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yangkequn/saavuu/api"
	_ "github.com/yangkequn/saavuu/https"
)

type Demo struct {
	Text string
}
type Demo1 struct {
	Text   string `mapstructure:"Text,nonempty"`
	Attach *Demo
}

var ApiDemo = api.Api(func(InParam *Demo1) (ret string, err error) {
	now := time.Now()
	fmt.Println("Demo api is called with InParam:" + InParam.Text + " run at " + now.String() + " Attach:" + InParam.Attach.Text)
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
	var (
		err   error
		now   time.Time = time.Now()
		param           = &Demo1{Text: "TestCallAt 10s later", Attach: &Demo{Text: "Attach"}}
	)

	fmt.Println("Demo api is calling with InParam:" + param.Text + " run at " + now.String())

	callAt := api.CallAt(DemoRpc, now.Add(time.Second*10))

	if err = callAt(param); err != nil {
		t.Error(err)
	}
	time.Sleep(15 * time.Second)
}
func TestCallAtCancel(t *testing.T) {
	var (
		err   error
		now   time.Time = time.Now()
		param           = &Demo1{Text: "TestCallAt 10s later"}
	)
	timeToRun := time.Now().Add(time.Second * 10)
	callAt := api.CallAt(DemoRpc, timeToRun)
	fmt.Println("Demo api is calling with InParam:" + param.Text + " run at " + now.String())

	if err = callAt(param); err != nil {
		t.Error(err)
	}
	if ok := api.CallAtCancel(DemoRpc, timeToRun); !ok {
		t.Error("cancel failed")
	}
	time.Sleep(30 * time.Second)
}
