package service

import (
	"fmt"
	"saavuu/http"
)

func init() {
	type Input struct {
		data []uint16
	}

	http.NewService("Svc:Demo", func(svcCtx *http.HttpContext) (data interface{}, err error) {
		JwtID, ok := svcCtx.JwtField("id").(string)
		if !ok {
			return nil, http.ErrJWT
		}
		var in = &Input{}
		if err = http.ToStruct(svcCtx.Req, in); err != nil {
			return nil, err
		}
		// your logic here
		fmt.Println("JwtID is ", JwtID)

		return data, nil
	})
}
