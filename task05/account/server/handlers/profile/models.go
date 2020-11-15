package profile

import (
	"github.com/ds-vologdin/otus-software-architect/task05/account/users"
)

type Profile struct {
	ID        users.UserID
	Username  string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// UserToProfile convert User to Profile
func UserToProfile(user users.User) Profile {
	profile := Profile{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}
	return profile
}
