package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ds-vologdin/otus-software-architect/task06/app/auth/config"
)

const ()

var (
	ErrPasswordInvalid = errors.New("password is invalid")
	ErrInternal        = errors.New("internal error")

	DefaultHTTPRequestTimeout = 1 * time.Second
)

type UserID string

func UserIDFromInt64(i int64) UserID {
	return UserID(fmt.Sprintf("%d", i))
}

func (id UserID) String() string {
	return string(id)
}

type AccountProvider interface {
	CheckPassword(username, password string) (UserID, error)
}

type accountProvider struct {
	URL string
}

func (p accountProvider) CheckPassword(username, password string) (UserID, error) {
	client := &http.Client{
		Timeout: DefaultHTTPRequestTimeout,
	}

	req, err := http.NewRequest("GET", p.URL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s:%s", username, password))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", ErrPasswordInvalid
	}
	if resp.StatusCode != http.StatusOK {
		return "", ErrInternal
	}

	var response = struct{ ID int64 }{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	return UserIDFromInt64(response.ID), nil
}

func NewAccountProvider(conf config.AccountServiceConfig) AccountProvider {
	return accountProvider{URL: conf.URL}
}
