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
	"github.com/yangkequn/saavuu/permission"

	"github.com/vmihailenco/msgpack/v5"
)

// listten to a port and start http server
func RedisHttpStart(path string, port int) {
	var (
		result      interface{}
		b           []byte
		s           string
		ok          bool
		isByteArray bool
		isString    bool
		err         error
	)
	//get item
	router := http.NewServeMux()
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if https.CorsChecked(r, w) {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		defer cancel()
		svcContext := https.NewHttpContext(ctx, r, w)
		if r.Method == "GET" {
			result, err = svcContext.GetHandler()
		} else if r.Method == "POST" {
			result, err = svcContext.PutHandler()
		} else if r.Method == "DELETE" {
			result, err = svcContext.DelHandler()
		}
		if len(config.Cfg.CORS) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", config.Cfg.CORS)
		}
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "JWT") {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write([]byte(errStr))
			return
		}
		if len(svcContext.ResponseContentType) > 0 && len(svcContext.QueryFields) > 0 {
			w.Header().Set("Content-Type", svcContext.ResponseContentType)
		}
		w.WriteHeader(http.StatusOK)
		if b, isByteArray = result.([]byte); isByteArray {
			w.Write(b)
		} else if s, isString = result.(string); isString {
			w.Write([]byte(s))
		} else {
			//keep fields exits in svcContext.QueryFields
			if len(svcContext.QueryFields) > 0 {
				_map := map[string]interface{}{}
				//check if type of result is not map[string]interface{}
				if _map, ok = result.(map[string]interface{}); !ok {
					//convert result to map[string]interface{} using msgpack
					if b, err = msgpack.Marshal(result); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
						return
					}
					if err = msgpack.Unmarshal(b, &_map); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
						return
					}
				}
				//remove fields not exits in svcContext.QueryFields
				if svcContext.QueryFields != "*" {
					for k := range _map {
						if !strings.Contains(svcContext.QueryFields, k) {
							delete(_map, k)
						}
					}
				}
				result = _map
				//reponse result json to client
				if b, err = json.Marshal(result); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
				w.Write(b)
			}
		}
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

func main() {

	logger.Std.Println("App Start! load config from OS env")
	config.LoadConfigFromEnv()
	go permission.RefreshDataItemBatchPermissions()

	RedisHttpStart("/rSvc", 8080)
}
