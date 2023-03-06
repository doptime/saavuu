package test

import (
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
)

var apiTest = api.New("test")

func CreateTestApi() {
	apiTest.Serve(func(parmIn map[string]interface{}) (data map[string]interface{}, err error) {
		// your logic here
		data = map[string]interface{}{"data": "ok"}
		fmt.Println("test api ok")
		return data, nil
	})
}
func TestApi() {
	now := time.Now()
	go apiTest.DoAt(map[string]string{"message": "ok"}, &now)
	api.RunningAllApis()
}
