package saavuu

import (
	"fmt"
	"net/http"
	"strconv"
)

// listten to a port and start http server
func ListenAndServe(cfg *Configuration) {
	//get item
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		svcContext := LoadHttpContext(r, w)
		if svcContext.Jwt == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if svcContext.QueryFields == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if result, err := svcContext.getHandler(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if len(result) > 0 {
				w.Write(result)
			}
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		}
	})
	http.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q	", r.URL.Path)
	})
	http.HandleFunc("/rm", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q	", r.URL.Path)
	})
	http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
}
