package api

import (
	"net/http"
	"strings"

	"github.com/benharmonics/backend/logging"
)

func cors(h http.Handler, allowMethods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for the preflight request
		if r.Method == http.MethodOptions {
			logging.Debug("CORS preflight request received")
			methods := strings.Join(allowMethods, " ")
			if len(allowMethods) == 0 {
				methods = "*"
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Mas-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Set CORS headers for the main request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
