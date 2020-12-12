package server

import (
	"net/http"

	"github.com/ds-vologdin/otus-software-architect/task06/app/billing/bill"

	"github.com/gorilla/mux"
)

type billServer struct {
	Billing bill.BillService
}

func (s *billServer) create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *billServer) delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *billServer) expense(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *billServer) getBalance(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *billServer) topUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// RegisterSubrouter register subrouter for work with a bill of user
func RegisterSubrouterBillig(base *mux.Router, path string, billService bill.BillService) error {
	s := billServer{Billing: billService}
	r := base.PathPrefix(path).Subrouter()
	r.HandleFunc("/{id}", s.create).Methods("POST")
	r.HandleFunc("/{id}", s.delete).Methods("DELETE")
	r.HandleFunc("/{id}/expense", s.expense).Methods("POST")
	r.HandleFunc("/{id}/balance", s.getBalance).Methods("GET")
	r.HandleFunc("/{id}/top_up", s.topUp).Methods("POST")
	return nil
}
