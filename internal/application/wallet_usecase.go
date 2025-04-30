package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/config"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/repository"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/cache"
	errs "github.com/hydr0g3nz/wallet_topup_system/internal/domain/error"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/logger"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/transaction"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/user"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/wallet"
)

type WalletUsecase interface {
	VerifyTopup(userID uint, amount float64, paymentMethod string) (transaction.Transaction, error)
	ConfirmTopup(transactionID uint) (transaction.Transaction, wallet.Wallet, error)
}

// WalletUsecase handles the business logic for wallet top-up operations
type WalletUsecaseImpl struct {
	userRepo        user.Repository
	transactionRepo transaction.Repository
	walletRepo      wallet.Repository
	cache           cache.CacheService
	tx              domain.DBTransaction // atomic transaction
	repoTx          domain.Repository
	logger          logger.Logger
	cfg             config.Config
}

// NewWalletUsecase creates a new instance of WalletUsecase
func NewWalletUsecase(
	userRepo user.Repository,
	transactionRepo transaction.Repository,
	walletRepo wallet.Repository,
	cache cache.CacheService,
	tx domain.DBTransaction,
	logger logger.Logger,
	config config.Config,

) WalletUsecase {
	repoTransaction := repository.NewRepositoryTransaction(transactionRepo, walletRepo, userRepo)
	return &WalletUsecaseImpl{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		cache:           cache,
		tx:              tx,
		repoTx:          repoTransaction,
		logger:          logger,
		cfg:             config,
	}
}

// VerifyTopup verifies a top-up request and creates a transaction with "verified" status
func (uc *WalletUsecaseImpl) VerifyTopup(userID uint, amount float64, paymentMethod string) (transaction.Transaction, error) {
	// Check if user exists
	if amount > uc.cfg.App.MaxAcceptedAmount {
		return transaction.Transaction{}, errs.ErrAmountExceedsLimit
	}
	_, err := uc.userRepo.FindById(userID)
	if err != nil {
		return transaction.Transaction{}, err
	}
	newTransaction, err := transaction.NewTransaction(userID, amount, paymentMethod, string(vo.StatusVerified), time.Now().Add(15*time.Minute))
	if err != nil {
		return transaction.Transaction{}, err
	}
	// Save transaction
	id, err := uc.transactionRepo.Create(newTransaction)
	if err != nil {
		return transaction.Transaction{}, err
	}
	newTransaction.ID = id

	// Store in cache (using transaction ID as key)
	cacheKey := getTransactionCacheKey(id)
	err = uc.cache.Set(context.Background(), cacheKey, newTransaction, 15*time.Minute)
	if err != nil {
		uc.logger.Error("Failed to set transaction in cache", map[string]interface{}{"error": err})
	}

	return newTransaction, nil
}

// ConfirmTopup confirms a previously verified transaction and updates the wallet balance
func (uc *WalletUsecaseImpl) ConfirmTopup(transactionID uint) (transaction.Transaction, wallet.Wallet, error) {
	// Try to get transaction from cache first
	cacheKey := getTransactionCacheKey(transactionID)
	tx := &transaction.Transaction{}
	err := uc.cache.Get(context.Background(), cacheKey, tx)
	if err != nil {
		uc.logger.Error("Failed to get transaction from cache", map[string]interface{}{"error": err})
		// Get transaction from database if not found in cache
		tx, err = uc.transactionRepo.FindById(transactionID)
		if err != nil {
			return transaction.Transaction{}, wallet.Wallet{}, err
		}
	} else {
		uc.logger.Info("Transaction found in cache", map[string]interface{}{"transaction": tx})
		uc.logger.Info("Transaction", map[string]interface{}{"transaction": tx})
	}

	// Check if transaction is verified and not expired
	if tx.Status != vo.StatusVerified {
		return transaction.Transaction{}, wallet.Wallet{}, errs.ErrTransactionNotVerified
	}

	if time.Now().After(tx.ExpiresAt) {
		// Update status to expired
		status := vo.StatusVerified
		err = uc.transactionRepo.Update(&transaction.TransactionFilter{ID: &tx.ID, Status: &status}, transaction.Transaction{
			Status: vo.StatusExpired,
		})
		if err != nil {
			return transaction.Transaction{}, wallet.Wallet{}, err
		}
		return transaction.Transaction{}, wallet.Wallet{}, errs.ErrExpiredTransaction
	}

	// Get user's wallet
	userWallet, err := uc.walletRepo.FindById(tx.UserID)
	if err != nil {
		return transaction.Transaction{}, wallet.Wallet{}, err
	}
	//update value
	userWallet.Balance = userWallet.Balance.Add(tx.Amount)
	tx.Status = vo.StatusCompleted
	// Update transaction status to completed
	err = uc.tx.DoInTransaction(func(repo domain.Repository) error {
		// Update wallet balance
		err = repo.WalletRepository().Update(*userWallet)
		if err != nil {
			return err
		}
		// Update transaction status to completed
		status := vo.StatusVerified
		err = repo.TransactionRepository().Update(&transaction.TransactionFilter{ID: &tx.ID, Status: &status}, transaction.Transaction{
			Status: tx.Status,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return transaction.Transaction{}, wallet.Wallet{}, err
	}
	// Remove from cache
	_ = uc.cache.Delete(context.Background(), cacheKey)

	return *tx, *userWallet, nil
}

// Helper function to generate cache key for transaction
func getTransactionCacheKey(transactionID uint) string {
	return "transaction:" + fmt.Sprintf("%d", transactionID)
}
