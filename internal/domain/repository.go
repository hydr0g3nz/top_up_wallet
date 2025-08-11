package domain

import (
	"context"

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

type TxManager interface {
	BeginTx(ctx context.Context) (context.Context, error) // Begin tx และ return ctx ใหม่ที่มี tx embed
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}
type txKey struct{}

func WithTx(ctx context.Context, tx any) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetTx(ctx context.Context) any {
	return ctx.Value(txKey{})
}
