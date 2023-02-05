package test

import (
	"context"
	"fmt"
	"time"

	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/data"
)

func TestApi() {

	api.NewApiService("test", func(dc *data.DataCtx, pc *api.ApiCtx, parmIn map[string]interface{}) (data map[string]interface{}, err error) {
		// your logic here
		data = map[string]interface{}{"data": "ok"}
		fmt.Println("test api ok")
		return data, nil
	})
	pc := api.NewApiContext(context.Background())
	go pc.Api("test", map[string]string{"message": "ok"}, nil, time.Now().UnixMicro()+1000)
	api.RunningAllApis()
}
