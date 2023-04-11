package test

import (
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
)

type testReq struct {
	Data string
}

var apiTest = api.Api("demo", func(parmIn *testReq) (data string, err error) {
	// your logic here
	fmt.Println("test api ok" + parmIn.Data)
	return "apiTestedSuccess", nil
})

func CreateTestApi() {

}
func TestApi() {
	now := time.Now()
	go apiTest.DoAt(&testReq{"messageok"}, &now)
	api.RunningAllApis()
}
