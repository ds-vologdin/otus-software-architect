package bill

// BillService - service for work with users entities
type BillService interface {
	GetBalance(UserID) (Balance, error)
	Create(UserID) error
	Expense(Transaction) (Balance, Transaction, error)
	TopUp(Transaction) (Balance, Transaction, error)
	Delete(UserID) error
}
