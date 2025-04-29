package transaction

type Repository interface {
	FindAll(filter *TransactionFilter) ([]Transaction, error)
	FindById(id uint) (*Transaction, error)
	Create(transaction Transaction) (uint, error)
	Update(transaction Transaction) error
}
