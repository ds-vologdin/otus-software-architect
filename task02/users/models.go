package users

import (
	"fmt"
	"strconv"
)

// UserID - surrogate key of the user table
type UserID int64

// User - model of user
type User struct {
	ID        UserID `db:"id"`
	Username  string `db:"username"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Phone     string `db:"phone"`
}

func UserIDFromString(s string) (UserID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid format of UserID (%s)", s)
	}
	return UserID(i), nil
}
