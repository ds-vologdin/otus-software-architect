package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ds-vologdin/otus-software-architect/task06/app/billing/bill"
)

type billServer struct {
	Billing bill.BillService
}

func (s *billServer) create(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CREATE] request: %v %v", r.Method, r.URL)

	userID, ok := getUserIDFormRequestContext(r)
	if !ok {
		log.Printf("user id is empty in context of Request")
		http.Error(w, "create user failed", http.StatusInternalServerError)
		return
	}

	err := s.Billing.Create(userID)
	if err != nil {
		log.Printf("create bill: %v", err)
		switch err {
		case bill.ErrUserExists:
			http.Error(w, "user already exists", http.StatusConflict)
		default:
			http.Error(w, "create user failed", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *billServer) delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELETE] request: %v %v", r.Method, r.URL)

	userID, ok := getUserIDFormRequestContext(r)
	if !ok {
		log.Printf("user id is empty in context of Request")
		http.Error(w, "delete user failed", http.StatusInternalServerError)
		return
	}

	err := s.Billing.Delete(userID)
	if err != nil {
		http.Error(w, "delete user failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *billServer) expense(w http.ResponseWriter, r *http.Request) {
	log.Printf("[EXPENSE] request: %v %v", r.Method, r.URL)

	userID, ok := getUserIDFormRequestContext(r)
	if !ok {
		log.Printf("user id is empty in context of Request")
		http.Error(w, "delete user failed", http.StatusInternalServerError)
		return
	}

	var req struct {
		Amount  int64 `json:"amount"`
		OrderID int64 `json:"order_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("invalid format of request: %v", err)
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	log.Printf("expense: %+v", req)

	transaction := bill.Transaction{
		UserID:  userID,
		OrderID: bill.OrderID(req.OrderID),
		Type:    bill.Expense,
		Amount:  -req.Amount,
	}

	balance, transaction, err := s.Billing.Expense(transaction)
	if err != nil {
		http.Error(w, "expense failed", http.StatusBadRequest)
		return
	}

	resp := struct {
		Balance     bill.Balance     `json:"balance"`
		Transaction bill.Transaction `json:"transaction"`
	}{
		Balance:     balance,
		Transaction: transaction,
	}

	buff, err := json.Marshal(resp)
	if err != nil {
		log.Printf("marshal error: %v (%v)", err, resp)
		http.Error(w, "expense failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(buff)
}

func (s *billServer) topUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *billServer) getBalance(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// RegisterSubrouter register subrouter for work with a bill of user
func RegisterSubrouterBilling(base *mux.Router, path string, billService bill.BillService) error {
	s := billServer{Billing: billService}
	r := base.PathPrefix(path).Subrouter()
	r.HandleFunc("/{id}", s.create).Methods("POST")
	r.HandleFunc("/{id}", s.delete).Methods("DELETE")
	r.HandleFunc("/{id}/expense", s.expense).Methods("POST")
	r.HandleFunc("/{id}/balance", s.getBalance).Methods("GET")
	r.HandleFunc("/{id}/top_up", s.topUp).Methods("POST")

	r.Use(checkAccessMiddleware)
	return nil
}
