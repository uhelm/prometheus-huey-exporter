package exporter

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MakeHTTPHandler(metricsPath string) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(metricsPath, promhttp.Handler())

	return mux
}
