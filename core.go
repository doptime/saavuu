package saavuu

import (
	"net/http"
	"strconv"
)

// listten to a port and start http server
func ListenAndServe(cfg *Configuration) {
	var (
		result []byte
		err    error
	)
	//get item
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		svcContext := LoadHttpContext(r, w)
		if r.Method == "GET" {
			result, err = svcContext.getHandler()
		} else if r.Method == "PUT" {
			result, err = svcContext.putHandler()
		} else if r.Method == "DELETE" {
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
			w.Write(result)
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
}
