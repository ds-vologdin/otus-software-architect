package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ds-vologdin/otus-software-architect/task05/account/users"
)

const (
	maxConnectionRetry = 5
)

type service struct {
	db *sqlx.DB
}

// NewUserService create new UserService
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
	err := svc.db.Get(&user, `
		SELECT id, username, first_name, last_name, email, phone
		FROM "user" WHERE id=$1;`,
		id,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return user, fmt.Errorf("get user %v: %w", id, users.UserNotFound)
	}
	return user, err
}

func (svc *service) CheckCredential(username, password string) (users.UserID, error) {
	user := struct {
		ID       users.UserID
		Password string
	}{}
	err := svc.db.Get(&user, `SELECT id, password FROM "user" WHERE username=$1;`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("check credential for %v: %w", username, users.UserNotFound)
		}
		log.Printf("check credential for %v: %v", username, err)
		return 0, users.ErrUserGet
	}
	if user.Password != HashString(password) {
		log.Printf("check credential for %v: invalid password", username)
		return 0, users.ErrInvalidPassword
	}
	return user.ID, nil
}

func (svc *service) Create(user users.User) (users.UserID, error) {
	password := HashString(user.Password)
	var id users.UserID

	tx, err := svc.db.Beginx()
	if err != nil {
		log.Printf("create user: %v", err)
		return id, users.ErrUserCreate
	}

	var exists bool
	err = tx.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM "user" WHERE username=$1);`,
		user.Username,
	)
	if exists {
		log.Printf("create user: user already exists")
		return id, users.ErrUserExists
	}

	err = tx.Get(&id, `
		INSERT INTO "user" (username, password, first_name, last_name, email, phone)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id;`,
		user.Username, password, user.FirstName, user.LastName, user.Email, user.Phone,
	)
	if err != nil {
		log.Printf("create user: %v", err)
		return id, users.ErrUserCreate
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("create user: %v", err)
		return id, users.ErrUserCreate
	}
	return id, nil
}

func (svc *service) Edit(user users.User) error {
	tx, err := svc.db.Beginx()
	if err != nil {
		log.Printf("edit user: %v", err)
		return users.ErrUserEdit
	}

	userFromDb := users.User{}
	err = tx.Get(&userFromDb, `
		SELECT id, first_name, last_name, email, phone
		FROM "user" WHERE id=$1
		FOR UPDATE;`,
		user.ID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get user %v: %w", user.ID, users.UserNotFound)
	}

	user = mergeUsers(user, userFromDb)

	r, err := tx.Exec(`
		UPDATE "user"
		SET first_name=$2, last_name=$3, email=$4, phone=$5
		WHERE id=$6;`,
		user.FirstName, user.LastName, user.Email, user.Phone, user.ID,
	)
	if err != nil {
		log.Printf("edit user: %v", err)
		return users.ErrUserEdit
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("edit user: %v", err)
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

func mergeUsers(user, userFromDB users.User) users.User {
	if user.Password != "" {
		user.Password = HashString(user.Password)
	} else {
		user.Password = userFromDB.Password
	}
	if user.FirstName == "" {
		user.FirstName = userFromDB.FirstName
	}
	if user.LastName == "" {
		user.LastName = userFromDB.LastName
	}
	if user.Email == "" {
		user.Email = userFromDB.Email
	}
	if user.Phone == "" {
		user.Phone = userFromDB.Phone
	}
	return user
}

func (svc *service) Update(user users.User) error {
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
		log.Printf("delete user: %v", err)
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
