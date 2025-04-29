package transaction

type Repository interface {
	FindAll() ([]Transaction, error)
	FindById(id uint) (Transaction, error)
	Save(transaction Transaction) error
}
