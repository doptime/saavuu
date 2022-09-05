package service

import (
	"fmt"
	"saavuu/https"
)

func init() {
	type Input struct {
		data []uint16
	}

	https.NewService("Svc:Demo", func(svcCtx *https.HttpContext) (data interface{}, err error) {
		JwtID, ok := svcCtx.JwtField("id").(string)
		if !ok {
			return nil, https.ErrJWT
		}
		var in = &Input{}
		if err = https.ToStruct(svcCtx.Req, in); err != nil {
			return nil, err
		}
		// your logic here
		fmt.Println("JwtID is ", JwtID)

		return data, nil
	})
}
