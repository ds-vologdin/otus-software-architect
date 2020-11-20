package users

import "errors"

var (
	UserNotFound       = errors.New("user is not found")
	ErrUserEdit        = errors.New("edit user has failed")
	ErrUserDelete      = errors.New("delete user has failed")
	ErrUserCreate      = errors.New("create user has failed")
	ErrUserGet         = errors.New("get user has failed")
	ErrInvalidPassword = errors.New("invalid password")
)
