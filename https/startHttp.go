package https

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/permission"
)

// listten to a port and start http server
func RedisHttpStart(path string, port int64) {
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
			svcCtx     *HttpContext
		)
		if CorsChecked(r, w) {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*12000)
		defer cancel()
		if svcCtx, err = NewHttpContext(ctx, r, w); err != nil {
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
		} else if b, err = json.Marshal(result); err != nil {
			//reponse result json to client
			httpStatus = http.StatusInternalServerError
			b = []byte(err.Error())
		}

		svcCtx.SetContentType()
		w.WriteHeader(httpStatus)
		w.Write(b)
	})
	log.Info().Any("port ", port).Any("path is ", path).Msg("http server started")

	server := &http.Server{
		Addr:              ":" + strconv.FormatInt(port, 10),
		Handler:           router,
		ReadTimeout:       50 * time.Second,
		ReadHeaderTimeout: 50 * time.Second,
		WriteTimeout:      50 * time.Second, //10ms Redundant time
		IdleTimeout:       15 * time.Second,
	}
	server.ListenAndServe()
}
func StartHttp() {
	permission.LoadPPermissionFromRedis()
	RedisHttpStart("/", config.Cfg.ServerPort)
}
