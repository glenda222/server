package healthcheck

import (
	"net/http"
	"strings"
)

func Heartbeat(endpoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == "GET" || r.Method == "HEAD") &&
			strings.EqualFold(r.URL.Path, endpoint) {
			// other checks are performed here
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return
		}
	}
	return http.HandlerFunc(fn)
}
