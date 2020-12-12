package bill

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ds-vologdin/otus-software-architect/task06/app/billing/config"
)

const (
	maxConnectionRetry = 5
)

type service struct {
	db *sqlx.DB
}

func NewBillService(cfg config.DatabaseConfig) (BillService, error) {
	var db *sqlx.DB
	var err error
	for i := 0; i < maxConnectionRetry; i++ {
		db, err = sqlx.Connect("postgres", cfg.DSN)
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

func (svc *service) GetBalance(id UserID) (Balance, error) {
	balance := Balance{}
	err := svc.db.Get(&balance, `
		SELECT *
		FROM balance WHERE user_id=$1;`,
		id,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return balance, ErrUserNotFound
	}
	if err != nil {
		log.Printf("get balance: %v", err)
		return balance, ErrGetBalance
	}
	return balance, nil
}

func (svc *service) Create(id UserID) error {
	tx, err := svc.db.Beginx()
	if err != nil {
		log.Printf("create bill of user: %v", err)
		return ErrCreateBill
	}

	var exists bool
	err = tx.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM balance WHERE id=$1);`,
		id,
	)
	if exists {
		log.Printf("create bill of user: user already exists")
		return ErrUserExists
	}

	err = tx.Get(&id, `
		INSERT INTO balance (user_id, amount)
		VALUES ($1, 0)
		RETURNING id;`,
		id,
	)
	if err != nil {
		log.Printf("create bill of user: %v", err)
		return ErrCreateBill
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("create bill of user: %v", err)
		return ErrCreateBill
	}
	return nil
}

func (svc *service) Expense(transaction Transaction) (Balance, Transaction, error) {
	balance := Balance{UserID: transaction.UserID}

	if transaction.Amount > 0 {
		return balance, transaction, fmt.Errorf("expense: amount > 0 (%d)", transaction.Amount)
	}
	if transaction.Type != Expense {
		return balance, transaction, fmt.Errorf("expense: type of transaction != %s (%s)", Expense, transaction.Type)
	}

	tx, err := svc.db.Beginx()
	if err != nil {
		log.Printf("expense: %v", err)
		return balance, transaction, ErrCreateTransaction
	}

	err = tx.Get(&balance.Amount, `
		SELECT amount FROM balance WHERE user_id=$1 FOR UPDATE;`,
		balance.UserID,
	)
	if err != nil {
		log.Printf("expense: %v", err)
		return balance, transaction, ErrCreateTransaction
	}

	if balance.Amount < transaction.Amount {
		transaction.Status = Declined
	} else {
		transaction.Status = Accepted
	}
	transaction.Time = time.Now().UTC()

	var id TransactionID
	err = tx.Get(&id, `
		INSERT INTO transaction (user_id, order_id, time, type, amount, status)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id;`,
		transaction.UserID, transaction.OrderID, transaction.Time, transaction.Amount, transaction.Status,
	)
	if err != nil {
		log.Printf("expense: %v", err)
		return balance, transaction, ErrCreateTransaction
	}

	if transaction.Status == Accepted {
		balance.Amount -= transaction.Amount
		_, err = tx.Exec(`
		UPDATE balance SET amount=$1 WHERE user_id=$2;`,
			balance.Amount, balance.UserID,
		)
		if err != nil {
			log.Printf("expense: %v", err)
			return balance, transaction, ErrCreateTransaction
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("expense: %v", err)
		return balance, transaction, ErrCreateTransaction
	}
	return balance, transaction, nil
}

func (svc *service) TopUp(transaction Transaction) (Balance, Transaction, error) {
	balance := Balance{UserID: transaction.UserID}

	if transaction.Amount < 0 {
		return balance, transaction, fmt.Errorf("top up: amount < 0 (%d)", transaction.Amount)
	}
	if transaction.Type != TopUp {
		return balance, transaction, fmt.Errorf("top up: type of transaction != %s (%s)", TopUp, transaction.Type)
	}

	tx, err := svc.db.Beginx()
	if err != nil {
		log.Printf("top up: %v", err)
		return balance, transaction, ErrCreateTransaction
	}

	err = tx.Get(&balance.Amount, `
		SELECT amount FROM balance WHERE user_id=$1 FOR UPDATE;`,
		balance.UserID,
	)
	if err != nil {
		log.Printf("top up: %v", err)
		return balance, transaction, ErrCreateTransaction
	}

	transaction.Status = Accepted
	transaction.Time = time.Now().UTC()

	err = tx.Get(&transaction.ID, `
		INSERT INTO transaction (user_id, order_id, time, type, amount, status)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id;`,
		transaction.UserID, transaction.OrderID, transaction.Time, transaction.Amount, transaction.Status,
	)
	if err != nil {
		log.Printf("top up: %v", err)
		return balance, transaction, ErrCreateTransaction
	}

	balance.Amount += transaction.Amount
	_, err = tx.Exec(`
		UPDATE balance SET amount=$1 WHERE user_id=$2;`,
		balance.Amount, balance.UserID,
	)
	if err != nil {
		log.Printf("top up: %v", err)
		return balance, transaction, ErrCreateTransaction
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("top up: %v", err)
		return balance, transaction, ErrCreateTransaction
	}
	return balance, transaction, nil
}

func (svc *service) Delete(id UserID) error {
	tx, err := svc.db.Beginx()
	if err != nil {
		log.Printf("top up: %v", err)
		return ErrDeleteBill
	}

	_, err = tx.Exec(`DELETE FROM transaction WHERE user_id=$1;`, id)
	if err != nil {
		log.Printf("delete: %v", err)
		return ErrDeleteBill
	}

	_, err = tx.Exec(`DELETE FROM balance WHERE user_id=$1;`, id)
	if err != nil {
		log.Printf("delete: %v", err)
		return ErrDeleteBill
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("delete: %v", err)
		return ErrDeleteBill
	}
	return nil
}
