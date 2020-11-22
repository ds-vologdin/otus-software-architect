package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ds-vologdin/otus-software-architect/task05/auth-sidecar/config"
)

const (
	maxShutdownTime = 10 * time.Second
)

// Server is struct with HTTP server and UserService
type Server struct {
	SVC *http.Server
}

// Handlers
func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(MsgStatusOK)
}

// End handlers

// NewServer - initialize the http server
func NewServer(cfg config.Config) (*Server, error) {
	s := Server{}

	r := mux.NewRouter()
	r.HandleFunc("/healthz", health)
	r.Handle("/metrics", promhttp.Handler())

	proxy, err := NewProxy(cfg)
	if err != nil {
		log.Printf("create proxy: %v", err)
		return nil, err
	}
	r.PathPrefix("/").HandlerFunc(proxy.Proxy)

	r.Use(metricsMiddleware, headerMiddleware)
	address := fmt.Sprintf(":%d", cfg.Server.Port)
	s.SVC = &http.Server{Addr: address, Handler: r}

	log.Printf("Run HTTP proxy: %v -> %v", address, proxy.Target)

	return &s, nil
}

// Run - function for run server. Support graceful shutdown.
func (srv *Server) Run() {
	shutdown := make(chan struct{})
	defer close(shutdown)
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stop
		log.Printf("Got signal '%v', the graceful shutdown will start", sig)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, maxShutdownTime)
		defer cancel()

		if err := srv.SVC.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		} else {
			log.Print("HTTP server has been shutdown")
		}
		shutdown <- struct{}{}
	}()

	if err := srv.SVC.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	log.Print("Wait for the shutdown server")
	<-shutdown
}
