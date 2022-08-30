package service

import (
	"saavuu/http"
)

func init() {
	type Input struct {
		data []uint16
	}

	http.NewService("Svc:Demo", func(svcCtx *http.HttpContext) (data interface{}, err error) {
		var in = &Input{}
		if err = http.ToStruct(svcCtx.Req, in); err != nil {
			return nil, err
		}
		// your logic here

		return data, nil
	})
}
