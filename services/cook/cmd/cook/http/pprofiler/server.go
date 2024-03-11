package pprofiler

import (
	"net/http"
	"net/http/pprof"
)

func NewServer(addr string) *http.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return &http.Server{
		Addr:    ":6060",
		Handler: handler,
	}
}
