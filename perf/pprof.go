package perf

import (
	"net/http"
	_ "net/http/pprof"
)

func ServeMonitor(addr string) error {
	return http.ListenAndServe(addr, nil)
}
