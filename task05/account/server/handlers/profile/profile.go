package profile

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ds-vologdin/otus-software-architect/task05/account/users"
	"github.com/gorilla/mux"
)

type server struct {
	UserService users.UserService
}

func (srv *server) getProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("[GET USER] request: %v %v [%v]", r.Method, r.URL, vars)

	xUserID := r.Header.Get("X-User-Id")
	if xUserID != vars["id"] {
		log.Printf("X-User-Id invalid: %v != %v", xUserID, vars["id"])
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

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

func (srv *server) createUser(w http.ResponseWriter, r *http.Request) {
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

func (srv *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("[DELETE USER] request: %v %v [%v]", r.Method, r.URL, vars)

	xUserID := r.Header.Get("X-User-Id")
	if xUserID != vars["id"] {
		log.Printf("X-User-Id invalid: %v != %v", xUserID, vars["id"])
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

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

func (srv *server) editProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("[EDIT USER] request: %v %v [%v]", r.Method, r.URL, vars)

	xUserID := r.Header.Get("X-User-Id")
	if xUserID != vars["id"] {
		log.Printf("X-User-Id invalid: %v != %v", xUserID, vars["id"])
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	id, err := users.UserIDFromString(vars["id"])
	if err != nil {
		log.Printf("invalid id: %v", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(MsgInvalidUserID)
		return
	}

	var user users.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("user format error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MsgInvalidDataFormat)
		return
	}

	if user.ID == 0 {
		user.ID = id
	} else {
		if user.ID != id {
			log.Printf("user id from json != user id from url (%v != %v)", user.ID, id)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(MsgInvalidUserID)
			return
		}
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

// RegisterSubrouter register subrouter for work with profile of user
func RegisterSubrouter(base *mux.Router, path string, userService users.UserService) error {
	s := server{userService}
	r := base.PathPrefix(path).Subrouter()
	r.HandleFunc("/", s.createUser).Methods("POST")
	r.HandleFunc("/{id}", s.getProfile).Methods("GET")
	r.HandleFunc("/{id}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/{id}", s.editProfile).Methods("PUT")
	return nil
}
