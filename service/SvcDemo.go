package service

import (
	"saavuu/http"
)

func init() {
	type DemoInput struct {
		data []uint16
		//ogg *multipart.FileHeader
	}

	http.NewService("SvcDemo", func(svcCtx *http.HttpContext) (data interface{}, err error) {
		var i = &DemoInput{}
		if err = http.ToStruct(svcCtx.Req, i); err != nil {
			return nil, err
		}
		// your logic here

		return data, nil
	})
}
