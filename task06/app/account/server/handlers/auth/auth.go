package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ds-vologdin/otus-software-architect/task06/app/account/users"
	"github.com/gorilla/mux"
)

var (
	errNeedAuth    = errors.New("need basic authorization")
	errInvalidAuth = errors.New("invalid format of basic authorization")

	msgCheckCredentialFailed = "check credential failed"
	msgAuthorizationFailed   = "authorization failed"
)

type Credential struct {
	Username string
	Password string
}

type server struct {
	UserService users.UserService
}

func (s *server) checkCredential(w http.ResponseWriter, r *http.Request) {
	cred, err := getCredentialsFromAuthorization(r)
	if err != nil {
		authorizeError(w)
		return
	}

	userID, err := s.UserService.CheckCredential(cred.Username, cred.Password)
	if err != nil {
		if errors.Is(err, users.ErrInvalidPassword) {
			authorizeError(w)
			return
		}
		log.Printf("check credential for %v: %v", cred.Username, err)
		http.Error(w, msgCheckCredentialFailed, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(struct{ ID users.UserID }{userID})
	if err != nil {
		log.Printf("encode to json error: %v", err)
		http.Error(w, msgCheckCredentialFailed, http.StatusInternalServerError)
		return
	}
}

// RegisterSubrouter register subrouter for work with profile of user
func RegisterSubrouter(base *mux.Router, path string, userService users.UserService) error {
	s := server{userService}
	r := base.PathPrefix(path).Subrouter()
	r.HandleFunc("", s.checkCredential).Methods("GET")
	return nil
}

func getCredentialsFromAuthorization(r *http.Request) (Credential, error) {
	var cred Credential

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		return cred, errNeedAuth
	}

	credString := auth[1]
	splited := strings.SplitN(credString, ":", 2)
	if len(splited) != 2 {
		return cred, errInvalidAuth
	}

	cred.Username = splited[0]
	cred.Password = splited[1]
	return cred, nil
}

func authorizeError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic")
	http.Error(w, msgAuthorizationFailed, http.StatusUnauthorized)
}
