package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"github.com/ds-vologdin/otus-software-architect/task05/auth/config"
	"github.com/ds-vologdin/otus-software-architect/task05/auth/providers/account"
)

const (
	TokenTypeRefresh = "refresh"
	TokenTypeAccess  = "access"

	TokenTypeRefreshTTL = time.Duration(24 * time.Hour)
	TokenTypeAccessTTL  = time.Duration(15 * time.Minute)
)

var (
	errNeedAuth    = errors.New("need basic authorization")
	errInvalidAuth = errors.New("invalid format of basic authorization")
)

type server struct {
	AccountProvider account.AccountProvider
	PrivateKey      *rsa.PrivateKey
	PublicKey       *rsa.PublicKey
}

type Credential struct {
	Username string
	Password string
}

type RefreshToken struct {
	RefreshToken string
	AccessToken  string
}

type ParsedToken struct {
	UID       string
	Type      string
	ExpiresOn int64
}

func (s *server) createRefreshToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("createRefreshToken")
	cred, err := getCredentialsFromAuthorization(r)
	if err != nil {
		authorizeError(w)
		return
	}
	log.Printf("cred: %v", cred)

	userID, err := s.AccountProvider.CheckPassword(cred.Username, cred.Password)
	if err != nil {
		if errors.Is(err, account.ErrPasswordInvalid) {
			authorizeError(w)
			return
		}
		log.Printf("check credential for %v: %v", cred.Username, err)
		http.Error(w, "check password error", http.StatusInternalServerError)
		return
	}
	log.Printf("user id: %v", userID)

	refreshToken, err := createToken(s.PrivateKey, userID.String(), TokenTypeRefresh, TokenTypeRefreshTTL)
	if err != nil {
		log.Printf("create refresh token error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	log.Printf("refreshToken: %v", refreshToken)

	accessToken, err := createToken(s.PrivateKey, userID.String(), TokenTypeAccess, TokenTypeAccessTTL)
	if err != nil {
		log.Printf("create access token error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	log.Printf("accessToken: %v", accessToken)

	tokens := RefreshToken{refreshToken, accessToken}
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(tokens)
	if err != nil {
		log.Printf("encode to json error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (s *server) createAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshTokens := r.Header["Authorization"]
	if refreshTokens == nil {
		log.Print("request without Authorization header field")
		w.Header().Set("WWW-Authenticate:", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(MsgRefreshTokenEmpty)
	}

	refreshTokenString := refreshTokens[0]
	if strings.HasPrefix(refreshTokenString, "Bearer ") {
		refreshTokenString = refreshTokenString[7:]
	}

	refreshToken, err := parseToken(s.PublicKey, refreshTokenString)
	if err != nil {
		log.Printf("invalid refresh token: %v", err)
		w.WriteHeader(http.StatusForbidden)
		w.Write(MsgInvalidRefreshToken)
	}
	if err = checkParsedRefreshToken(refreshToken); err != nil {
		log.Printf("invalid refresh token: %v", err)
		w.WriteHeader(http.StatusForbidden)
		w.Write(MsgInvalidRefreshToken)
	}

	accessToken, err := newAccessToken(s.PrivateKey, refreshToken)
	if err != nil {
		log.Printf("create access token error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MsgInternalError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(struct{ AccessToken string }{accessToken})
	if err != nil {
		log.Printf("encode to json error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MsgInternalError)
		return
	}
}

// RegisterSubrouter register subrouter for work with token of user
func RegisterSubrouter(base *mux.Router, path string, accountProvider account.AccountProvider, cfg config.JWTConfig) error {
	if cfg.Algorithm != "RS256" {
		return fmt.Errorf("unsupport jwt algorithm '%s', support only RS256", cfg.Algorithm)
	}
	privateKeyRaw, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	if err != nil {
		log.Printf("invalid jwt private key: %v", err)
		return err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyRaw)
	if err != nil {
		log.Printf("invalid jwt private key: %v", err)
		return err
	}

	if isInvalidPublicKey(privateKey, cfg.PublicKey) {
		log.Printf("invalid jwt public key")
		return errors.New("invalid public key")
	}

	s := server{
		AccountProvider: accountProvider,
		PrivateKey:      privateKey,
		PublicKey:       &privateKey.PublicKey,
	}
	r := base.PathPrefix(path).Subrouter()
	r.HandleFunc("/refresh", s.createRefreshToken).Methods("POST")
	r.HandleFunc("/access", s.createAccessToken).Methods("POST")
	return nil
}

func createToken(privateKey *rsa.PrivateKey, uid, tokenType string, period time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"uid":  uid,
		"type": tokenType,
		"exp":  time.Now().Add(period).Unix(),
	})
	return token.SignedString(privateKey)
}

func newAccessToken(privateKey *rsa.PrivateKey, refreshToken ParsedToken) (string, error) {
	return createToken(privateKey, refreshToken.UID, TokenTypeAccess, TokenTypeAccessTTL)
}

func parseToken(publicKey *rsa.PublicKey, tokenString string) (ParsedToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	parsedToken := ParsedToken{}

	if err != nil {
		log.Printf("invalid token: %v", token)
		return parsedToken, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Printf("refresh token: %v", token)
		return parsedToken, fmt.Errorf("invalid token")
	}

	v := claims["uid"]
	if uid, ok := v.(string); ok {
		parsedToken.UID = uid
	}
	v = claims["type"]
	if tokenType, ok := v.(string); ok {
		parsedToken.Type = tokenType
	}
	v = claims["exp"]
	if exp, ok := v.(int64); ok {
		parsedToken.ExpiresOn = exp
	}

	return parsedToken, nil
}

func checkParsedRefreshToken(token ParsedToken) error {
	if token.Type != TokenTypeRefresh {
		return fmt.Errorf("invalid type of token: want '%v', got '%v'", TokenTypeRefresh, token.Type)
	}
	if token.UID == "" {
		return fmt.Errorf("token has not uid")
	}
	if token.ExpiresOn > time.Now().Unix() {
		return fmt.Errorf("token is expired")
	}
	return nil
}

func authorizeError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic")
	http.Error(w, "authorization failed", http.StatusUnauthorized)
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

func isInvalidPublicKey(privateKey *rsa.PrivateKey, gotPublicKey string) bool {
	publicKeyRaw := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicKey := base64.StdEncoding.EncodeToString(publicKeyRaw)
	return publicKey != gotPublicKey
}
