package test

import (
	"context"
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
)

func CreateTestApi() {
	api.NewApi("test", func(pc *api.Ctx, parmIn map[string]interface{}) (data map[string]interface{}, err error) {
		// your logic here
		data = map[string]interface{}{"data": "ok"}
		fmt.Println("test api ok")
		return data, nil
	})
}
func TestApi() {
	_api := api.NewContext(context.Background())
	now := time.Now()
	go _api.DoAt("test", map[string]string{"message": "ok"}, &now)
	api.RunningAllApis()
}
