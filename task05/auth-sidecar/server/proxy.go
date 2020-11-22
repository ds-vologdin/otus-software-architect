package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ds-vologdin/otus-software-architect/task05/auth-sidecar/config"
)

const (
	TokenTypeAccess    = "access"
	TokenTypeAccessTTL = time.Duration(15 * time.Minute)
)

var (
	DefaultHTTPRequestTimeout = 1 * time.Second

	ErrUnauthorized = errors.New("unauthorized request")
	ErrForbidden    = errors.New("forbidden request")
)

type proxyServer struct {
	PublicKey   *rsa.PublicKey
	Target      string
	ExcludeAuth []config.Request
}

type ParsedToken struct {
	UID       string
	Type      string
	ExpiresOn int64
}

func (p proxyServer) Proxy(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, " ", r.Method, " ", r.URL)

	var userID string
	if isNeedAuth(p.ExcludeAuth, r) {
		parsedToken, err := auth(p.PublicKey, r)
		if err != nil {
			switch err {
			case ErrUnauthorized:
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			case ErrForbidden:
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}
		userID = parsedToken.UID
	}

	client := &http.Client{
		Timeout: DefaultHTTPRequestTimeout,
	}

	url := fmt.Sprintf("%s%s", p.Target, r.URL.Path)
	log.Printf("%v %v", r.Method, url)
	req, err := http.NewRequest(r.Method, url, r.Body)

	copyHeader(req.Header, r.Header)
	if userID != "" {
		req.Header.Add("X-User-Id", userID)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("request: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func auth(key *rsa.PublicKey, r *http.Request) (ParsedToken, error) {
	var parsedToken ParsedToken
	token := r.Header.Get("Authorization")
	if token == "" {
		log.Print("request without Authorization header field")
		return parsedToken, ErrUnauthorized
	}

	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	log.Printf("token: %v", token)

	parsedToken, err := parseToken(key, token)
	if err != nil {
		log.Print(err)
		return parsedToken, ErrForbidden
	}
	log.Printf("parsed: %v", parsedToken)
	if err = checkParsedAccessToken(parsedToken); err != nil {
		log.Printf("invalid token: %v", err)
		return parsedToken, ErrForbidden
	}
	return parsedToken, nil
}

func isNeedAuth(exclude []config.Request, r *http.Request) bool {
	for _, request := range exclude {
		if request.Path != r.URL.Path {
			continue
		}
		if request.Method == r.Method {
			return false
		}
	}
	return true
}

func NewProxy(cfg config.Config) (proxyServer, error) {
	p := proxyServer{}
	if cfg.JWT.Algorithm != "RS256" {
		return p, fmt.Errorf("unsupport jwt algorithm '%s', support only RS256", cfg.JWT.Algorithm)
	}
	publicKeyRaw, err := base64.StdEncoding.DecodeString(cfg.JWT.PublicKey)
	if err != nil {
		log.Printf("invalid jwt public key: %v", err)
		return p, err
	}
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyRaw)
	if err != nil {
		log.Printf("invalid jwt public key: %v", err)
		return p, err
	}
	p.PublicKey = publicKey

	target := cfg.Target.URL
	targetLength := len(target)
	if strings.HasSuffix(target, "/") {
		target = target[:targetLength-1]
	}
	p.Target = target
	p.ExcludeAuth = cfg.ExcludeAuth
	return p, nil
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

func checkParsedAccessToken(token ParsedToken) error {
	if token.Type != TokenTypeAccess {
		return fmt.Errorf("invalid type of token: want '%v', got '%v'", TokenTypeAccess, token.Type)
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
	w.Header().Set("WWW-Authenticate", "Bearer")
	http.Error(w, "authorization failed", http.StatusUnauthorized)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
