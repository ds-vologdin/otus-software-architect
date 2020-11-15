package users

import "errors"

var (
	UserNotFound  = errors.New("user not found")
	ErrUserEdit   = errors.New("error user edit")
	ErrUserDelete = errors.New("erroru ser delete")
	ErrUserCreate = errors.New("error user create")
)
