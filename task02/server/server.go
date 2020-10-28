package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/ds-vologdin/otus-software-architect/task02/users"
)

const (
	maxShutdownTime = 10 * time.Second
)

// Server is struct with HTTP server and UserService
type Server struct {
	SVC         *http.Server
	UserService users.UserService
}

// Handlers
func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(MsgStatusOK)
}

func (srv *Server) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("[GET USER] request: %v %v [%v]", r.Method, r.URL, vars)

	id, err := users.UserIDFromString(vars["id"])
	if err != nil {
		log.Printf("invalid id: %v", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(MsgInvalidUserID)
		return
	}

	user, err := srv.UserService.Get(users.UserID(id))
	if err != nil {
		log.Printf("get user error: %v", err)
		w.WriteHeader(http.StatusNotFound)

		if errors.Is(err, users.UserNotFound) {
			w.Write(MsgUserNotFound)
		} else {
			w.Write(MsgGetUserError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MsgInternalError)
		return
	}
}

func (srv *Server) postUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[POST USER] request: %v %v", r.Method, r.URL)

	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("user format error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MsgInvalidDataFormat)
		return
	}

	id, err := srv.UserService.Create(user)
	if err != nil {
		log.Printf("create user error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MsgCreateUserError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(struct{ ID users.UserID }{id})
	if err != nil {
		log.Printf("encode to json error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MsgInternalError)
		return
	}
}

func (srv *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("[DELETE USER] request: %v %v [%v]", r.Method, r.URL, vars)

	id, err := users.UserIDFromString(vars["id"])
	if err != nil {
		log.Printf("invalid id: %v", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(MsgInvalidUserID)
		return
	}
	err = srv.UserService.Delete(id)
	if err != nil {
		log.Printf("delete user error: %v", err)
		if errors.Is(err, users.UserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(MsgUserNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(MsgDeleteUserError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(MsgStatusOK)
}

func (srv *Server) editUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("[EDIT USER] request: %v %v [%v]", r.Method, r.URL, vars)
	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("user format error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MsgInvalidDataFormat)
		return
	}

	err = srv.UserService.Edit(user)
	if err != nil {
		log.Printf("edit user error: %v", err)
		if errors.Is(err, users.UserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(MsgUserNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(MsgEditUserError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(MsgStatusOK)
}

// End handlers

// NewServer - initialize the http server
func NewServer(address string, userService users.UserService) (*Server, error) {
	s := Server{}
	s.UserService = userService

	r := mux.NewRouter()
	r.HandleFunc("/healthz", health)
	r.HandleFunc("/user/", s.postUser).Methods("POST")
	r.HandleFunc("/user/{id}", s.getUser).Methods("GET")
	r.HandleFunc("/user/{id}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/user/{id}", s.editUser).Methods("PUT")
	r.Use(headerMiddleware)
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
