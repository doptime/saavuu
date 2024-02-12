package test

import (
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
)

func TestApi() {

	type InTest struct {
		Data string
	}

	var apiTest = api.Api(func(parmIn *InTest) (data string, err error) {
		// your logic here
		fmt.Println("test api ok" + parmIn.Data)
		return "apiTestedSuccess", nil
	})
	apiTest(&InTest{"messageok"})

	var apiTestLaterDo = api.CallAt(time.Now().Add(10*time.Second), apiTest)
	go apiTestLaterDo(&InTest{"messageok"})
	api.StarApis()
}
