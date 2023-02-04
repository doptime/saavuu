package test

import (
	"context"
	"fmt"
	"time"

	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/rCtx"
)

func TestRpc() {

	saavuu.NewRpcService("test", func(dc *rCtx.DataCtx, pc *rCtx.ParamCtx, parmIn map[string]interface{}) (data map[string]interface{}, err error) {
		// your logic here
		data = map[string]interface{}{"data": "ok"}
		fmt.Println("test api ok")
		return data, nil
	})
	pc := saavuu.NewParamContext(context.Background())
	go pc.RpcBasic("test", map[string]string{"message": "ok"}, time.Now().UnixMicro()+1000)
	saavuu.RunningAllRpcs()
}
