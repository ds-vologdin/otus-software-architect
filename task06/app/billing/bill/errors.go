package bill

import "errors"

var (
	ErrUserNotFound      = errors.New("user is not found")
	ErrGetBalance        = errors.New("get a balance has failed")
	ErrCreateBill        = errors.New("create a bill has failed")
	ErrCreateTransaction = errors.New("create a transaction has failed")
	ErrDeleteBill        = errors.New("delete a bill has failed")
	ErrUserExists        = errors.New("user already exists")
)
