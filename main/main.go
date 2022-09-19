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
	"github.com/yangkequn/saavuu/tools"

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
				for _, k := range tools.MapKeys(_map) {
					if !strings.Contains(svcContext.QueryFields, k) {
						delete(_map, k)
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
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func main() {
	config.LoadConfigFromEnv("config/configsaavuu.toml")

	RedisHttpStart("/rSvc", 3025)
}
