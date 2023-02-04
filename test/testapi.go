package test

import (
	"context"
	"fmt"
	"time"

	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/rCtx"
)

func TestApi() {

	saavuu.NewService("test", func(dc *rCtx.DataCtx, pc *rCtx.ParamCtx, parmIn map[string]interface{}) (data map[string]interface{}, err error) {
		// your logic here
		data = map[string]interface{}{"data": "ok"}
		fmt.Println("test api ok")
		return data, nil
	})
	pc := saavuu.NewParamContext(context.Background())
	go pc.RdsApiBasic("test", map[string]string{"message": "ok"}, time.Second*10)
	saavuu.RunningAllService()
}
