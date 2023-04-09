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
			svcCtx     *https.HttpContext
		)
		if https.CorsChecked(r, w) {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*12000)
		defer cancel()
		if svcCtx, err = https.NewHttpContext(ctx, r, w); err != nil {
			httpStatus = http.StatusBadRequest
		} else if r.Method == "GET" {
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
			} else if httpStatus == http.StatusOK {
				httpStatus = http.StatusInternalServerError
			}
			b = []byte(err.Error())
		} else if b, ok = result.([]byte); ok {
		} else if s, ok = result.(string); ok {
			b = []byte(s)
		} else if len(svcCtx.ResponseFields) > 0 {
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
	RedisHttpStart("/", 8080)
}
