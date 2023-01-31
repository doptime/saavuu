package test

import (
	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/rCtx"
)

func testApi() {

	saavuu.NewService("test", 128, func(dc *rCtx.DataCtx, pc *rCtx.ParamCtx, parmIn map[string]interface{}) (data map[string]interface{}, err error) {
		// your logic here
		data = map[string]interface{}{"data": "ok"}
		return data, nil
	})
	// pc := saavuu.NewParamContext(context.Background())
	// go pc.RdsApiBasic("test", "message")
	saavuu.RunningAllService()
}
