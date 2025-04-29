package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/cache"
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
	logger          Logger
}

// NewWalletUsecase creates a new instance of WalletUsecase
func NewWalletUsecase(
	userRepo user.Repository,
	transactionRepo transaction.Repository,
	walletRepo wallet.Repository,
	cache cache.CacheService,
	tx domain.DBTransaction,
	logger Logger,

) WalletUsecase {
	return &WalletUsecaseImpl{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		cache:           cache,
		tx:              tx,
		logger:          logger,
	}
}

// VerifyTopup verifies a top-up request and creates a transaction with "verified" status
func (uc *WalletUsecaseImpl) VerifyTopup(userID uint, amount float64, paymentMethod string) (transaction.Transaction, error) {
	// Check if user exists
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
		uc.logger.Error("Failed to set transaction in cache", err)
	}

	return newTransaction, nil
}

// ConfirmTopup confirms a previously verified transaction and updates the wallet balance
func (uc *WalletUsecaseImpl) ConfirmTopup(transactionID uint) (transaction.Transaction, wallet.Wallet, error) {
	// Try to get transaction from cache first
	cacheKey := getTransactionCacheKey(transactionID)
	tx := &transaction.Transaction{}
	// var x transaction.Transaction
	err := uc.cache.Get(context.Background(), cacheKey, tx)
	if err != nil {
		uc.logger.Error("Failed to get transaction from cache", err)
		// Get transaction from database if not found in cache
		tx, err = uc.transactionRepo.FindById(transactionID)
		if err != nil {
			return transaction.Transaction{}, wallet.Wallet{}, err
		}
	} else {
		uc.logger.Info("Transaction found in cache", cacheKey)
		uc.logger.Info("Transaction found in cache", tx)
	}

	// Check if transaction is verified and not expired
	if tx.Status != vo.StatusVerified {
		return transaction.Transaction{}, wallet.Wallet{}, errors.New("transaction is not in 'verified' status")
	}

	if time.Now().After(tx.ExpiresAt) {
		// Update status to expired
		tx.Status = vo.StatusExpired
		err = uc.transactionRepo.Update(*tx)
		if err != nil {
			return transaction.Transaction{}, wallet.Wallet{}, err
		}
		return transaction.Transaction{}, wallet.Wallet{}, errors.New("transaction has expired")
	}

	// Get user's wallet
	userWallet, err := uc.walletRepo.FindById(tx.UserID)
	if err != nil {
		return transaction.Transaction{}, wallet.Wallet{}, err
	}
	// Update wallet balance
	userWallet.Balance = userWallet.Balance.Add(tx.Amount)
	err = uc.walletRepo.Update(*userWallet)
	if err != nil {
		return transaction.Transaction{}, wallet.Wallet{}, err
	}

	// Update transaction status to completed
	tx.Status = vo.StatusCompleted
	err = uc.transactionRepo.Update(transaction.Transaction{
		Status: tx.Status,
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
