package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ds-vologdin/otus-software-architect/task02/users"
)

const (
	maxConnectionRetry = 5
)

type service struct {
	db *sqlx.DB
}

func NewUserService(dsn string) (users.UserService, error) {
	var db *sqlx.DB
	var err error
	for i := 0; i < maxConnectionRetry; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		log.Printf("connect to DB error: %v", err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return &service{db}, nil
}

func (svc *service) Get(id users.UserID) (users.User, error) {
	user := users.User{}
	err := svc.db.Get(&user, `SELECT * FROM "user" WHERE id=$1;`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return user, fmt.Errorf("get user %v: %w", id, users.UserNotFound)
	}
	return user, err
}

func (svc *service) Create(user users.User) (users.UserID, error) {
	var id users.UserID
	sql := `
		INSERT INTO "user" (username, first_name, last_name, email, phone)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id;
	`
	err := svc.db.Get(&id, sql, user.Username, user.FirstName, user.LastName, user.Email, user.Phone)
	if err != nil {
		return id, users.ErrUserCreate
	}
	return id, nil
}

func (svc *service) Edit(user users.User) error {
	sql := `
		UPDATE "user"
		SET username=$1, first_name=$2, last_name=$3, email=$4, phone=$5
		WHERE id=$6;
	`
	r, err := svc.db.Exec(sql, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.ID)
	if err != nil {
		return users.ErrUserEdit
	}
	count, err := r.RowsAffected()
	if err != nil {
		return users.ErrUserEdit
	}
	if count == 0 {
		return users.UserNotFound
	}
	return nil
}

func (svc *service) Delete(id users.UserID) error {
	sql := `
		DELETE FROM "user" WHERE id=$1;
	`
	r, err := svc.db.Exec(sql, id)
	if err != nil {
		return users.ErrUserEdit
	}
	count, err := r.RowsAffected()
	if err != nil {
		return users.ErrUserDelete
	}
	if count == 0 {
		return users.UserNotFound
	}
	return nil
}
