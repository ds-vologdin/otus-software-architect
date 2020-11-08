package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestLatency *prometheus.HistogramVec
	requestCount   *prometheus.CounterVec
)

func init() {
	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "accounts_request_latency_sec",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)
	prometheus.MustRegister(requestLatency)

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "accounts_request_count",
		},
		[]string{"method", "endpoint", "http_status"},
	)
	prometheus.MustRegister(requestCount)
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		endpoint := parseEndpoint(r.URL.Path)

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(recorder, r)
		status := strconv.Itoa(recorder.Status)

		requestLatency.WithLabelValues(r.Method, endpoint).Observe(time.Since(start).Seconds())
		requestCount.WithLabelValues(r.Method, endpoint, status).Inc()
	})
}

func parseEndpoint(path string) string {
	if strings.HasPrefix(path, "/user/") {
		return "/user/"
	}
	return path
}
