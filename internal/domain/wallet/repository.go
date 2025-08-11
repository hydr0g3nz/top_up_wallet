package wallet

import "context"

type Repository interface {
	Create(wallet Wallet) error
	Update(ctx context.Context, wallet Wallet) error
	FindById(id uint) (*Wallet, error)
}
