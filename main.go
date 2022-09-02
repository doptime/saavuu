package main

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"
	. "saavuu/config"
	sttp "saavuu/http"
	"saavuu/service"
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
		svcContext := sttp.NewHttpContext(r, w)
		if r.Method == "GET" {
			result, err = svcContext.GetHandler()
		} else if r.Method == "POST" {
			result, err = svcContext.PutHandler()
		} else if r.Method == "DELETE" {
			result, err = svcContext.DelHandler()
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
		w.WriteHeader(http.StatusOK)
		if len(svcContext.ExpectedReponseType) > 0 {
			w.Header().Set("Content-Type", svcContext.ExpectedReponseType)
		}
		if b, ok = result.([]byte); ok {
			w.Write(b)
		} else if s, ok = result.(string); ok {
			w.Write([]byte(s))
		} else {
			//keep fields exits in svcContext.QueryFields
			if len(svcContext.QueryFields) > 0 {
				//convert result to map[string]interface{} using msgpack
				_tmpMap := map[string]interface{}{}
				if b, err = msgpack.Marshal(result); err != nil {
					sttp.InternalError(w, err)
					return
				}
				if err = msgpack.Unmarshal(b, &_tmpMap); err != nil {
					sttp.InternalError(w, err)
					return
				}
				//remove fields not exits in svcContext.QueryFields
				for k, _ := range _tmpMap {
					if !strings.Contains(svcContext.QueryFields, k) {
						delete(_tmpMap, k)
					}
				}
				//write to client
				if b, err = json.Marshal(_tmpMap); err != nil {
					sttp.InternalError(w, err)
					return
				}
			} else {
				if b, err = json.Marshal(result); err != nil {
					sttp.InternalError(w, err)
					return
				}
			}
			w.Write(b)
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
