package transaction

type Repository interface {
	FindAll(filter *TransactionFilter) ([]Transaction, error)
	FindById(id uint) (*Transaction, error)
	Create(transaction Transaction) (uint, error)
	Update(filter *TransactionFilter, transaction Transaction) error
}
