package saavuu

import (
	"net/http"
	"strconv"
)

// listten to a port and start http server
func RedisHttpStart(cfg *Configuration, path string) {
	var (
		result []byte
		err    error
	)
	//get item
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		svcContext := LoadHttpContext(r, w)
		if svcContext.Action == "GET" {
			result, err = svcContext.getHandler()
		} else if svcContext.Action == "PUT" {
			result, err = svcContext.putHandler()
		} else if svcContext.Action == "DELETE" {
			result, err = svcContext.delHandler()
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if len(result) > 0 {
				w.Write([]byte(err.Error()))
			}
			return
		} else {
			w.WriteHeader(http.StatusOK)
			if len(svcContext.ExpectedReponseType) > 0 {
				w.Header().Set("Content-Type", svcContext.ExpectedReponseType)
			}
			w.Write(result)
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
}
