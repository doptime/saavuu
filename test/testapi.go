package test

import (
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
)

type testReq struct {
	Data string
}

var apiTest = api.New[*testReq]("test")

func CreateTestApi() {
	apiTest.Serve(func(parmIn *testReq) (data interface{}, err error) {
		// your logic here
		fmt.Println("test api ok" + parmIn.Data)
		data = map[string]interface{}{"data": "ok"}
		return data, nil
	})
}
func TestApi() {
	now := time.Now()
	go apiTest.DoAt(map[string]string{"message": "ok"}, &now)
	api.RunningAllApis()
}
