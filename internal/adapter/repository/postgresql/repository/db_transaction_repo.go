package repository

import (
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/transaction"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/user"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/wallet"
	"gorm.io/gorm"
)

type DBTransactionRepository struct {
	db *gorm.DB
}

func NewDBTransactionRepository(db *gorm.DB) *DBTransactionRepository {
	return &DBTransactionRepository{db: db}
}
func (d *DBTransactionRepository) DoInTransaction(fn func(repo domain.Repository) error) error {
	// เริ่ม transaction
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	repoCtx := &RepositoryTransaction{
		transactionRepo: NewTransactionRepository(tx),
		walletRepo:      NewWalletRepository(tx),
		userRepo:        NewUserRepository(tx),
	}
	// ทำงานภายใน transaction
	if err := fn(repoCtx); err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction
	return tx.Commit().Error
}

type RepositoryTransaction struct {
	transactionRepo transaction.Repository
	walletRepo      wallet.Repository
	userRepo        user.Repository
}

func NewRepositoryTransaction(
	transactionRepo transaction.Repository,
	walletRepo wallet.Repository,
	userRepo user.Repository,
) *RepositoryTransaction {
	return &RepositoryTransaction{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		userRepo:        userRepo,
	}
}
func (r *RepositoryTransaction) UserRepository() user.Repository {
	return r.userRepo
}
func (r *RepositoryTransaction) WalletRepository() wallet.Repository {
	return r.walletRepo
}
func (r *RepositoryTransaction) TransactionRepository() transaction.Repository {
	return r.transactionRepo
}
