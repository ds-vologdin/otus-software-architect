package users

// UserService - service for work with users entities
type UserService interface {
	Get(UserID) (User, error)
	Create(User) (UserID, error)
	Edit(User) error
	Delete(UserID) error
}
