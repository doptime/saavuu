package test

import (
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
)

func TestApi() {
	now := time.Now()

	type InTest struct {
		Data string
	}

	var apiTest, apiTestCtx = api.Api(func(parmIn *InTest) (data string, err error) {
		// your logic here
		fmt.Println("test api ok" + parmIn.Data)
		return "apiTestedSuccess", nil
	})

	apiTest(&InTest{"messageok"})
	go apiTestCtx.DoAt(&InTest{"messageok"}, &now)
	api.StarApis()
}
