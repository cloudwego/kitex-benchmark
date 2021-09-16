package perf

import (
	"net/http"
	_ "net/http/pprof"
)

func ServeMonitor(addr string) error {
	mux := http.NewServeMux()
	return http.ListenAndServe(addr, mux)
}
