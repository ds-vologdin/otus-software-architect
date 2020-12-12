package bill

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// Transaction types
	Expense = "expense"
	TopUp   = "top_up"

	// Statuses
	Accepted = "accepted"
	Declined = "declined"
)

type UserID int64
type OrderID int64
type TransactionID int64
type TransactionType string

type Transaction struct {
	ID      TransactionID   `db:"id" json:"id"`
	UserID  UserID          `db:"user_id" json:"user_id"`
	OrderID OrderID         `db:"order_id" json:"order_id"`
	Time    time.Time       `db:"time" json:"created"`
	Type    TransactionType `db:"type" json:"type"`
	Amount  int64           `db:"amount" json:"amount"`
	Status  string          `db:"status" json:"status"`
}

type Balance struct {
	UserID UserID `db:"user_id"`
	Amount int64  `db:"amount"`
}

// UserIDFromString convert string to UserID type
func UserIDFromString(s string) (UserID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid format of UserID (%s)", s)
	}
	return UserID(i), nil
}
