package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ds-vologdin/otus-software-architect/task05/auth/providers/account"
	"github.com/ds-vologdin/otus-software-architect/task05/auth/server/handlers/token"
)

const (
	maxShutdownTime = 10 * time.Second
)

// Server is struct with HTTP server and UserService
type Server struct {
	SVC             *http.Server
	AccountProvider account.AccountProvider
}

// Handlers
func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(MsgStatusOK)
}

// End handlers

// NewServer - initialize the http server
func NewServer(address string, accountProvider account.AccountProvider) (*Server, error) {
	s := Server{}
	s.AccountProvider = accountProvider

	r := mux.NewRouter()
	r.HandleFunc("/healthz", health)
	r.Handle("/metrics", promhttp.Handler())

	token.RegisterSubrouter(r, "/token", accountProvider)

	r.Use(metricsMiddleware, headerMiddleware)
	s.SVC = &http.Server{Addr: address, Handler: r}

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
