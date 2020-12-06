package users

// UserService - service for work with users entities
type UserService interface {
	Get(UserID) (User, error)
	CheckCredential(string, string) (UserID, error)
	Create(User) (UserID, error)
	Edit(User) error
	Update(User) error
	Delete(UserID) error
}
