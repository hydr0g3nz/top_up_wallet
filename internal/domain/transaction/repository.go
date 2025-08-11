package transaction

import "context"

type Repository interface {
	FindAll(filter *TransactionFilter) ([]Transaction, error)
	FindById(id uint) (*Transaction, error)
	Create(ctx context.Context, transaction Transaction) (uint, error)
	Update(ctx context.Context, filter *TransactionFilter, transaction Transaction) error
}
