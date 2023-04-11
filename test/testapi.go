package test

import (
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
)

type InDemo struct {
	Data string
}

var apiTest = api.Api(func(parmIn *InDemo) (data string, err error) {
	// your logic here
	fmt.Println("test api ok" + parmIn.Data)
	return "apiTestedSuccess", nil
})

func CreateTestApi() {

}
func TestApi() {
	now := time.Now()
	go apiTest.DoAt(&InDemo{"messageok"}, &now)
	api.RunningAllApis()
}
