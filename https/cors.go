package https

import (
	"net/http"
	"strconv"

	"github.com/yangkequn/saavuu/config"
)

func CorsChecked(r *http.Request, w http.ResponseWriter) bool {
	if r.Method == "OPTIONS" && len(config.Cfg.Http.CORES) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Accept-Language, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Origin", config.Cfg.Http.CORES)
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(30*86400))
		w.Header().Set("Content-Type", "text/html; charset=ascii")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("CORS OK"))
		return true
	}
	return false
}
