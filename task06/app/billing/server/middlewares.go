package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ds-vologdin/otus-software-architect/task06/app/billing/bill"
	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	headerUserID = "X-User-Id"
)

var (
	errAccessForbidden = errors.New("access forbidden")
	errInvalidUserID   = errors.New("invalid user id")
)

var (
	requestLatency *prometheus.HistogramVec
	requestCount   *prometheus.CounterVec
)

func init() {
	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_request_latency_sec",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)
	prometheus.MustRegister(requestLatency)

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_request_count",
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

func checkAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := checkUserIDFromPath(r)
		if err != nil {
			log.Printf("check access: %v", err)
			switch err {
			case errAccessForbidden:
				http.Error(w, "access forbidden", http.StatusForbidden)
			case errInvalidUserID:
				http.Error(w, "invalid user id", http.StatusBadRequest)
			default:
				http.Error(w, "internal error", http.StatusBadRequest)
			}
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkUserIDFromPath(r *http.Request) (bill.UserID, error) {
	var userID bill.UserID
	vars := mux.Vars(r)
	xUserID := r.Header.Get(headerUserID)

	if xUserID == "" {
		return userID, errAccessForbidden
	}
	if xUserID != vars["id"] {
		return userID, errAccessForbidden
	}

	userID, err := bill.UserIDFromString(xUserID)
	if err != nil {
		return userID, errInvalidUserID
	}
	return userID, nil
}

func getUserIDFormRequestContext(r *http.Request) (bill.UserID, bool) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(bill.UserID)
	return userID, ok
}
