package main

import (
	"net/http"

	"github.com/grafana/pyroscope-go"
)

func main() {
	pyroscope.Start(pyroscope.Config{
		ApplicationName: "demo.cook",
		ServerAddress:   "http://localhost:4040",
	})

	http.ListenAndServe(":8080", nil)
}
