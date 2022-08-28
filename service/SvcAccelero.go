package service

import (
	"fmt"
	"saavuu/http"
)

func init() {
	type Accelero struct {
		Accleration1s []int16
	}

	http.NewService("svc:Acceleros1s", func(svcCtx *http.HttpContext) (data interface{}, err error) {
		var (
			i = &Accelero{}
		)
		if err = http.ToStruct(svcCtx.Req, i); err != nil || len(i.Accleration1s)%3 != 0 {
			return nil, err
		}

		for v := range i.Accleration1s {
			var a float32 = 16.0 * float32(v) / 32768.0
			fmt.Println("a:", a)
		}

		return data, nil
	})
}
