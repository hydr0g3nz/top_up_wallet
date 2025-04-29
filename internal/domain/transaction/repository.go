package transaction

type Repository interface {
	FindAll(filter TransactionFilter) ([]Transaction, error)
	FindById(id uint) (Transaction, error)
	Save(transaction Transaction) error
	Update(transaction Transaction) error
}
