package users

// UserService - service for work with users entities
type UserService interface {
	Get(UserID) (User, error)
	CheckPassword(UserID, string) (bool, error)
	Create(User) (UserID, error)
	Edit(User) error
	Delete(UserID) error
}
