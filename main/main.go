package main

import (
	"context"
	"time"

	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/https"
	"github.com/yangkequn/saavuu/logger"
)

// listten to a port and start http server
func RedisHttpStart(path string, port int) {
	//get item
	router := http.NewServeMux()
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var (
			result     interface{}
			b          []byte
			s          string
			ok         bool
			err        error
			httpStatus int = http.StatusOK
		)
		if https.CorsChecked(r, w) {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*12000)
		defer cancel()
		svcCtx := https.NewHttpContext(ctx, r, w)
		if r.Method == "GET" {
			result, err = svcCtx.GetHandler()
		} else if r.Method == "POST" {
			result, err = svcCtx.PostHandler()
		} else if r.Method == "PUT" {
			result, err = svcCtx.PutHandler()
		} else if r.Method == "DELETE" {
			result, err = svcCtx.DelHandler()
		}
		if len(config.Cfg.CORS) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", config.Cfg.CORS)
		}

		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "JWT") {
				httpStatus = http.StatusUnauthorized
			} else {
				httpStatus = http.StatusInternalServerError
			}
			b = []byte(err.Error())
		} else if b, ok = result.([]byte); ok {
		} else if s, ok = result.(string); ok {
			b = []byte(s)
		} else if len(svcCtx.QueryFields) > 0 {
			//keep fields exits in svcContext.QueryFields only
			if b, err = json.Marshal(result); err != nil {
				//reponse result json to client
				httpStatus = http.StatusInternalServerError
				b = []byte(err.Error())
			}
		}

		svcCtx.SetContentType()
		w.WriteHeader(httpStatus)
		w.Write(b)
	})
	logger.Lshortfile.Println("http server started on port " + strconv.Itoa(port) + " , path is " + path)

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           router,
		ReadTimeout:       50 * time.Second,
		ReadHeaderTimeout: 50 * time.Second,
		WriteTimeout:      50 * time.Second, //10ms Redundant time
		IdleTimeout:       15 * time.Second,
	}
	server.ListenAndServe()
}

type TestApi struct {
	ApiBase string
}

func main() {

	// db := data.NewContext(nil)
	// var info map[uint32]interface{}
	// if err := db.HGetAll("MeditBGInfo", &info); err != nil {
	// 	logger.Std.Println(err)
	// }
	// var info2 map[uint32]interface{} = make(map[uint32]interface{})
	// //take first 2 keys from info
	// keys2 := make([]uint32, 2)
	// i := 0
	// for k := range info {
	// 	keys2[i] = k
	// 	i++
	// 	if i == 2 {
	// 		break
	// 	}
	// }
	// if err := db.HMGET("MeditBGInfo", keys2, &info2); err != nil {
	// 	logger.Std.Println(err)
	// }
	// //print bytes
	// logger.Std.Println(len(info))

	//test.TestApi()

	RedisHttpStart("/rSvc", 8080)
}
