package domain

import (
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/transaction"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/user"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/wallet"
)

type Repository interface {
	UserRepository() user.Repository
	WalletRepository() wallet.Repository
	TransactionRepository() transaction.Repository
}
type DBTransaction interface {
	DoInTransaction(fn func(repo Repository) error) error
}
