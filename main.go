package main

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"
	. "saavuu/config"
	"saavuu/https"
	"saavuu/service"
	"saavuu/tools"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

// listten to a port and start http server
func RedisHttpStart(path string, port int) {
	var (
		result interface{}
		b      []byte
		s      string
		ok     bool
		err    error
	)
	//get item
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if https.CorsChecked(r, w) {
			return
		}
		svcContext := https.NewHttpContext(r, w)
		if r.Method == "GET" {
			result, err = svcContext.GetHandler()
		} else if r.Method == "POST" {
			result, err = svcContext.PutHandler()
		} else if r.Method == "DELETE" {
			result, err = svcContext.DelHandler()
		}
		if len(Cfg.CORS) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", Cfg.CORS)
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
		if len(svcContext.ExpectedReponseType) > 0 && len(svcContext.QueryFields) > 0 {
			w.Header().Set("Content-Type", svcContext.ExpectedReponseType)
		}
		w.WriteHeader(http.StatusOK)
		if b, ok = result.([]byte); ok {
			w.Write(b)
		} else if s, ok = result.(string); ok {
			w.Write([]byte(s))
		} else {
			//keep fields exits in svcContext.QueryFields
			if len(svcContext.QueryFields) > 0 {
				_map := map[string]interface{}{}
				//check if type of result is not map[string]interface{}
				if _map, ok = result.(map[string]interface{}); !ok {
					//convert result to map[string]interface{} using msgpack
					if b, err = msgpack.Marshal(result); err != nil {
						https.InternalError(w, err)
						return
					}
					if err = msgpack.Unmarshal(b, &_map); err != nil {
						https.InternalError(w, err)
						return
					}
				}
				//remove fields not exits in svcContext.QueryFields
				for _, k := range tools.MapKeys(_map) {
					if !strings.Contains(svcContext.QueryFields, k) {
						delete(_map, k)
					}
				}
				result = _map
				//reponse result json to client
				if b, err = json.Marshal(result); err != nil {
					https.InternalError(w, err)
					return
				}
				w.Write(b)
			}
		}
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func main() {
	Cfg.Rds = redis.NewClient(&redis.Options{
		Addr:     "docker.vm:6379", // use default Addr
		Password: "",               // no password set
		DB:       10,               // use default DB
	})
	Cfg.JwtSecret = Cfg.Rds.Get(context.Background(), "JwtSecret").String()
	fmt.Println("JwtSecret:", Cfg.JwtSecret)
	service.PrintServices()
	RedisHttpStart("/rSvc", 3025)
}
