package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/config"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain"
	errs "github.com/hydr0g3nz/wallet_topup_system/internal/domain/error"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/logger"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/transaction"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/user"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories and services
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(filter *user.UserFilter) ([]user.User, error) {
	args := m.Called(filter)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockUserRepository) FindById(id uint) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepository) Create(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) FindAll(filter *transaction.TransactionFilter) ([]transaction.Transaction, error) {
	args := m.Called(filter)
	return args.Get(0).([]transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindById(id uint) (*transaction.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Create(tx transaction.Transaction) (uint, error) {
	args := m.Called(tx)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockTransactionRepository) Update(filter *transaction.TransactionFilter, tx transaction.Transaction) error {
	args := m.Called(filter, tx)
	return args.Error(0)
}

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(w wallet.Wallet) error {
	args := m.Called(w)
	return args.Error(0)
}

func (m *MockWalletRepository) Update(w wallet.Wallet) error {
	args := m.Called(w)
	return args.Error(0)
}

func (m *MockWalletRepository) FindById(id uint) (*wallet.Wallet, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wallet.Wallet), args.Error(1)
}

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	// Simulate filling the destination if needed
	if len(args) > 2 && args.Get(2) != nil {
		switch v := dest.(type) {
		case *transaction.Transaction:
			tx := args.Get(2).(*transaction.Transaction)
			*v = *tx
		}
	}
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type MockDBTransaction struct {
	mock.Mock
}

func (m *MockDBTransaction) DoInTransaction(fn func(repo domain.Repository) error) error {
	args := m.Called(fn)
	// Execute the function with a mock repository if needed
	if len(args) > 1 && args.Get(1) != nil {
		err := fn(args.Get(1).(domain.Repository))
		if err != nil {
			return err
		}
	}
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Info(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Fatal(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) With(fields map[string]interface{}) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

type MockRepository struct {
	mock.Mock
	userRepo        *MockUserRepository
	walletRepo      *MockWalletRepository
	transactionRepo *MockTransactionRepository
}

func NewMockRepository(userRepo *MockUserRepository, walletRepo *MockWalletRepository, transactionRepo *MockTransactionRepository) *MockRepository {
	return &MockRepository{
		userRepo:        userRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (m *MockRepository) UserRepository() user.Repository {
	return m.userRepo
}

func (m *MockRepository) WalletRepository() wallet.Repository {
	return m.walletRepo
}

func (m *MockRepository) TransactionRepository() transaction.Repository {
	return m.transactionRepo
}

// Test cases
func TestVerifyTopup(t *testing.T) {
	// Setup mocks
	userRepo := new(MockUserRepository)
	transactionRepo := new(MockTransactionRepository)
	walletRepo := new(MockWalletRepository)
	cacheService := new(MockCacheService)
	dbTx := new(MockDBTransaction)
	logger := new(MockLogger)
	cfg := config.Config{
		App: config.AppConfig{
			MaxAcceptedAmount: 10000.0,
		},
	}

	// Create usecase with mocks
	walletUsecase := NewWalletUsecase(
		userRepo,
		transactionRepo,
		walletRepo,
		cacheService,
		dbTx,
		logger,
		cfg,
	)

	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		userID := uint(1)
		amount := 100.0
		paymentMethod := "credit_card"

		// Mock user existence check
		userRepo.On("FindById", userID).Return(user.User{ID: userID}, nil).Once()

		// Mock transaction creation
		transactionRepo.On("Create", mock.AnythingOfType("transaction.Transaction")).Return(uint(1), nil).Once()

		// Mock cache set
		cacheService.On("Set", mock.Anything, "transaction:1", mock.AnythingOfType("transaction.Transaction"), 15*time.Minute).Return(nil).Once()

		// Mock logger for potential errors
		logger.On("Error", mock.Anything, mock.Anything).Maybe()

		// Execute
		tx, err := walletUsecase.VerifyTopup(userID, amount, paymentMethod)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), tx.ID)
		assert.Equal(t, userID, tx.UserID)
		assert.Equal(t, vo.StatusVerified, tx.Status)

		// Verify mocks
		userRepo.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
		cacheService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup expectations
		userID := uint(999)
		amount := 100.0
		paymentMethod := "credit_card"

		// Mock user not found
		userRepo.On("FindById", userID).Return(user.User{}, errs.ErrNotFound).Once()

		// Execute
		_, err := walletUsecase.VerifyTopup(userID, amount, paymentMethod)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errs.ErrNotFound, err)

		// Verify mocks
		userRepo.AssertExpectations(t)
	})

	t.Run("Amount exceeds limit", func(t *testing.T) {
		// Setup expectations
		userID := uint(1)
		amount := 20000.0 // Greater than cfg.App.MaxAcceptedAmount (10000.0)
		paymentMethod := "credit_card"

		// Execute
		_, err := walletUsecase.VerifyTopup(userID, amount, paymentMethod)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errs.ErrAmountExceedsLimit, err)
	})

	t.Run("Invalid payment method", func(t *testing.T) {
		// Setup expectations
		userID := uint(1)
		amount := 100.0
		paymentMethod := "invalid_method"

		// Mock user existence check
		userRepo.On("FindById", userID).Return(user.User{ID: userID}, nil).Once()

		// Execute
		_, err := walletUsecase.VerifyTopup(userID, amount, paymentMethod)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidPaymentMethod, err)

		// Verify mocks
		userRepo.AssertExpectations(t)
	})

	t.Run("Transaction creation failure", func(t *testing.T) {
		// Setup expectations
		userID := uint(1)
		amount := 100.0
		paymentMethod := "credit_card"

		// Mock user existence check
		userRepo.On("FindById", userID).Return(user.User{ID: userID}, nil).Once()

		// Mock transaction creation failure
		transactionRepo.On("Create", mock.AnythingOfType("transaction.Transaction")).Return(uint(0), errors.New("database error")).Once()

		// Execute
		_, err := walletUsecase.VerifyTopup(userID, amount, paymentMethod)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		// Verify mocks
		userRepo.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
	})
}

func TestConfirmTopup(t *testing.T) {
	// Setup mocks
	userRepo := new(MockUserRepository)
	transactionRepo := new(MockTransactionRepository)
	walletRepo := new(MockWalletRepository)
	cacheService := new(MockCacheService)
	dbTx := new(MockDBTransaction)
	logger := new(MockLogger)
	cfg := config.Config{
		App: config.AppConfig{
			MaxAcceptedAmount: 10000.0,
		},
	}

	// Create usecase with mocks
	walletUsecase := NewWalletUsecase(
		userRepo,
		transactionRepo,
		walletRepo,
		cacheService,
		dbTx,
		logger,
		cfg,
	)

	// t.Run("Success with cache hit", func(t *testing.T) {
	// 	// Setup expectations
	// 	transactionID := uint(1)
	// 	userID := uint(1)

	// 	amount, _ := vo.NewMoney(100.0)
	// 	paymentMethod, _ := vo.NewPaymentMethod("credit_card")
	// 	status, _ := vo.NewTransactionStatus("verified")
	// 	expiresAt := time.Now().Add(10 * time.Minute) // Not expired

	// 	cachedTx := &transaction.Transaction{
	// 		ID:            transactionID,
	// 		UserID:        userID,
	// 		Amount:        amount,
	// 		PaymentMethod: paymentMethod,
	// 		Status:        status,
	// 		ExpiresAt:     expiresAt,
	// 	}

	// 	userWallet := &wallet.Wallet{
	// 		ID:      userID,
	// 		Balance: amount, // Initial balance equal to amount
	// 	}

	// 	expectedWallet := &wallet.Wallet{
	// 		ID:      userID,
	// 		Balance: amount.Add(amount), // Double the amount after top-up
	// 	}

	// 	// Mock cache get (hit)
	// 	cacheService.On("Get", mock.Anything, "transaction:1", mock.AnythingOfType("*transaction.Transaction")).
	// 		Run(func(args mock.Arguments) {
	// 			tx := args.Get(2).(*transaction.Transaction)
	// 			*tx = *cachedTx
	// 		}).
	// 		Return(nil).Once()

	// 	// Mock logger info
	// 	logger.On("Info", mock.Anything, mock.Anything).Return().Maybe()

	// 	// Mock wallet retrieval
	// 	walletRepo.On("FindById", userID).Return(userWallet, nil).Once()

	// 	// Mock transaction
	// 	mockRepo := NewMockRepository(userRepo, walletRepo, transactionRepo)
	// 	mockRepo.On("DoInTransaction", mock.AnythingOfType("func(domain.Repository) error")).
	// 		Run(func(args mock.Arguments) {
	// 			fn := args.Get(0).(func(interface{}) error)
	// 			fn(mockRepo) // Execute the function
	// 		}).
	// 		Return(nil).Once()

	// 	// Mock wallet update in transaction
	// 	walletRepo.On("Update", mock.MatchedBy(func(w wallet.Wallet) bool {
	// 		return w.ID == userID && w.Balance.Amount() == 200.0
	// 	})).Return(nil).Once()

	// 	// Mock transaction update in transaction
	// 	transactionRepo.On("Update",
	// 		mock.MatchedBy(func(filter *transaction.TransactionFilter) bool {
	// 			return *filter.ID == transactionID && *filter.Status == vo.StatusVerified
	// 		}),
	// 		mock.MatchedBy(func(tx transaction.Transaction) bool {
	// 			return tx.Status == vo.StatusCompleted
	// 		}),
	// 	).Return(nil).Once()

	// 	// Mock cache delete
	// 	cacheService.On("Delete", mock.Anything, "transaction:1").Return(nil).Once()

	// 	// Execute
	// 	resultTx, resultWallet, err := walletUsecase.ConfirmTopup(transactionID)

	// 	// Assert
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, vo.StatusCompleted, resultTx.Status)
	// 	assert.Equal(t, expectedWallet.Balance.Amount(), resultWallet.Balance.Amount())

	// 	// Verify mocks
	// 	cacheService.AssertExpectations(t)
	// 	walletRepo.AssertExpectations(t)
	// 	transactionRepo.AssertExpectations(t)
	// 	mockRepo.AssertExpectations(t)
	// })
	t.Run("Success with cache hit", func(t *testing.T) {
		// Setup expectations
		transactionID := uint(1)
		userID := uint(1)

		amount, _ := vo.NewMoney(100.0)
		paymentMethod, _ := vo.NewPaymentMethod("credit_card")
		status, _ := vo.NewTransactionStatus("verified")
		expiresAt := time.Now().Add(10 * time.Minute)

		cachedTx := &transaction.Transaction{
			ID:            transactionID,
			UserID:        userID,
			Amount:        amount,
			PaymentMethod: paymentMethod,
			Status:        status,
			ExpiresAt:     expiresAt,
		}

		userWallet := &wallet.Wallet{
			ID:      userID,
			Balance: amount,
		}

		// Mock logger to accept any calls - ข้ามการตรวจสอบ logger ทั้งหมด
		logger.On("Info", mock.Anything, mock.Anything).Return()
		logger.On("Error", mock.Anything, mock.Anything).Return()

		// Mock cache hit
		cacheService.On("Get", mock.Anything, "transaction:1", mock.AnythingOfType("*transaction.Transaction")).
			Run(func(args mock.Arguments) {
				tx := args.Get(2).(*transaction.Transaction)
				*tx = *cachedTx
			}).
			Return(nil).Once()

		// Mock wallet retrieval
		walletRepo.On("FindById", userID).Return(userWallet, nil).Once()

		// Mock DB transaction
		dbTx.On("DoInTransaction", mock.AnythingOfType("func(domain.Repository) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(0).(func(domain.Repository) error)
				err := fn(NewMockRepository(userRepo, walletRepo, transactionRepo))
				assert.NoError(t, err)
			}).
			Return(nil).Once()

		// Mock wallet update
		walletRepo.On("Update", mock.MatchedBy(func(w wallet.Wallet) bool {
			return w.ID == userID && w.Balance.Amount() == 200.0
		})).Return(nil).Once()

		// Mock transaction update
		transactionRepo.On("Update",
			mock.MatchedBy(func(filter *transaction.TransactionFilter) bool {
				return *filter.ID == transactionID && *filter.Status == vo.StatusVerified
			}),
			mock.MatchedBy(func(tx transaction.Transaction) bool {
				return tx.Status == vo.StatusCompleted
			}),
		).Return(nil).Once()

		// Mock cache deletion
		cacheService.On("Delete", mock.Anything, "transaction:1").Return(nil).Once()

		// Execute
		resultTx, resultWallet, err := walletUsecase.ConfirmTopup(transactionID)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, vo.StatusCompleted, resultTx.Status)
		assert.Equal(t, 200.0, resultWallet.Balance.Amount())

		// Verify mocks (ไม่ต้องตรวจสอบ logger)
		cacheService.AssertExpectations(t)
		walletRepo.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
		dbTx.AssertExpectations(t)
	})
	t.Run("Success with cache miss", func(t *testing.T) {
		// Setup expectations
		transactionID := uint(2)
		userID := uint(1)

		amount, _ := vo.NewMoney(100.0)
		paymentMethod, _ := vo.NewPaymentMethod("credit_card")
		status, _ := vo.NewTransactionStatus("verified")
		expiresAt := time.Now().Add(10 * time.Minute) // Not expired

		dbTxData := &transaction.Transaction{
			ID:            transactionID,
			UserID:        userID,
			Amount:        amount,
			PaymentMethod: paymentMethod,
			Status:        status,
			ExpiresAt:     expiresAt,
		}

		userWallet := &wallet.Wallet{
			ID:      userID,
			Balance: amount, // Initial balance equal to amount
		}

		// Mock cache get (miss)
		cacheService.On("Get", mock.Anything, "transaction:2", mock.AnythingOfType("*transaction.Transaction")).
			Return(errors.New("cache miss")).Once()

		// Mock logger error for cache miss
		logger.On("Error", mock.Anything, mock.Anything).Return().Maybe()

		// Mock transaction retrieval from database
		transactionRepo.On("FindById", transactionID).Return(dbTxData, nil).Once()

		// Mock wallet retrieval
		walletRepo.On("FindById", userID).Return(userWallet, nil).Once()

		// Mock transaction
		dbTx.On("DoInTransaction", mock.AnythingOfType("func(domain.Repository) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(0).(func(domain.Repository) error)
				err := fn(NewMockRepository(userRepo, walletRepo, transactionRepo))
				assert.NoError(t, err)
			}).
			Return(nil).Once()
		// Mock wallet update in transaction
		walletRepo.On("Update", mock.MatchedBy(func(w wallet.Wallet) bool {
			return w.ID == userID && w.Balance.Amount() == 200.0
		})).Return(nil).Once()

		// Mock transaction update in transaction
		transactionRepo.On("Update",
			mock.MatchedBy(func(filter *transaction.TransactionFilter) bool {
				return *filter.ID == transactionID && *filter.Status == vo.StatusVerified
			}),
			mock.MatchedBy(func(tx transaction.Transaction) bool {
				return tx.Status == vo.StatusCompleted
			}),
		).Return(nil).Once()

		// Mock cache delete
		cacheService.On("Delete", mock.Anything, "transaction:2").Return(nil).Once()

		// Execute
		resultTx, resultWallet, err := walletUsecase.ConfirmTopup(transactionID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, vo.StatusCompleted, resultTx.Status)
		assert.Equal(t, 200.0, resultWallet.Balance.Amount())

		// Verify mocks
		cacheService.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
		walletRepo.AssertExpectations(t)
		dbTx.AssertExpectations(t)
	})

	t.Run("Transaction not found", func(t *testing.T) {
		// Setup expectations
		transactionID := uint(999)

		// Mock cache get (miss)
		cacheService.On("Get", mock.Anything, "transaction:999", mock.AnythingOfType("*transaction.Transaction")).
			Return(errors.New("cache miss")).Once()

		// Mock logger error for cache miss
		logger.On("Error", mock.Anything, mock.Anything).Return().Maybe()

		// Mock transaction not found in database
		transactionRepo.On("FindById", transactionID).Return(nil, errs.ErrNotFound).Once()

		// Execute
		_, _, err := walletUsecase.ConfirmTopup(transactionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errs.ErrNotFound, err)

		// Verify mocks
		cacheService.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
	})

	t.Run("Transaction not in verified status", func(t *testing.T) {
		// Setup expectations
		transactionID := uint(3)
		userID := uint(1)

		amount, _ := vo.NewMoney(100.0)
		paymentMethod, _ := vo.NewPaymentMethod("credit_card")
		status, _ := vo.NewTransactionStatus("completed") // Already completed
		expiresAt := time.Now().Add(10 * time.Minute)

		cachedTx := &transaction.Transaction{
			ID:            transactionID,
			UserID:        userID,
			Amount:        amount,
			PaymentMethod: paymentMethod,
			Status:        status,
			ExpiresAt:     expiresAt,
		}

		// Mock cache get (hit)
		cacheService.On("Get", mock.Anything, "transaction:3", mock.AnythingOfType("*transaction.Transaction")).
			Run(func(args mock.Arguments) {
				tx := args.Get(2).(*transaction.Transaction)
				*tx = *cachedTx
			}).
			Return(nil).Once()

		// Mock logger info
		logger.On("Info", mock.Anything, mock.Anything).Return().Maybe()

		// Execute
		_, _, err := walletUsecase.ConfirmTopup(transactionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errs.ErrTransactionNotVerified, err)

		// Verify mocks
		cacheService.AssertExpectations(t)
	})

	t.Run("Transaction expired", func(t *testing.T) {
		// Setup expectations
		transactionID := uint(4)
		userID := uint(1)

		amount, _ := vo.NewMoney(100.0)
		paymentMethod, _ := vo.NewPaymentMethod("credit_card")
		status, _ := vo.NewTransactionStatus("verified")
		expiresAt := time.Now().Add(-10 * time.Minute) // Expired

		cachedTx := &transaction.Transaction{
			ID:            transactionID,
			UserID:        userID,
			Amount:        amount,
			PaymentMethod: paymentMethod,
			Status:        status,
			ExpiresAt:     expiresAt,
		}

		// Mock cache get (hit)
		cacheService.On("Get", mock.Anything, "transaction:4", mock.AnythingOfType("*transaction.Transaction")).
			Run(func(args mock.Arguments) {
				tx := args.Get(2).(*transaction.Transaction)
				*tx = *cachedTx
			}).
			Return(nil).Once()

		// Mock logger info
		logger.On("Info", mock.Anything, mock.Anything).Return().Maybe()

		// Mock transaction update to expired
		transactionRepo.On("Update",
			mock.MatchedBy(func(filter *transaction.TransactionFilter) bool {
				return *filter.ID == transactionID && *filter.Status == vo.StatusVerified
			}),
			mock.MatchedBy(func(tx transaction.Transaction) bool {
				return tx.Status == vo.StatusExpired
			}),
		).Return(nil).Once()

		// Execute
		_, _, err := walletUsecase.ConfirmTopup(transactionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errs.ErrExpiredTransaction, err)

		// Verify mocks
		cacheService.AssertExpectations(t)
		transactionRepo.AssertExpectations(t)
	})

	t.Run("Wallet not found", func(t *testing.T) {
		// Setup expectations
		transactionID := uint(5)
		userID := uint(999) // Non-existent wallet

		amount, _ := vo.NewMoney(100.0)
		paymentMethod, _ := vo.NewPaymentMethod("credit_card")
		status, _ := vo.NewTransactionStatus("verified")
		expiresAt := time.Now().Add(10 * time.Minute) // Not expired

		cachedTx := &transaction.Transaction{
			ID:            transactionID,
			UserID:        userID,
			Amount:        amount,
			PaymentMethod: paymentMethod,
			Status:        status,
			ExpiresAt:     expiresAt,
		}

		// Mock cache get (hit)
		cacheService.On("Get", mock.Anything, "transaction:5", mock.AnythingOfType("*transaction.Transaction")).
			Run(func(args mock.Arguments) {
				tx := args.Get(2).(*transaction.Transaction)
				*tx = *cachedTx
			}).
			Return(nil).Once()

		// Mock logger info
		logger.On("Info", mock.Anything, mock.Anything).Return().Maybe()

		// Mock wallet not found
		walletRepo.On("FindById", userID).Return(nil, errs.ErrNotFound).Once()

		// Execute
		_, _, err := walletUsecase.ConfirmTopup(transactionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errs.ErrNotFound, err)

		// Verify mocks
		cacheService.AssertExpectations(t)
		walletRepo.AssertExpectations(t)
	})

	t.Run("Transaction error in database operation", func(t *testing.T) {
		// Setup expectations
		transactionID := uint(6)
		userID := uint(1)

		amount, _ := vo.NewMoney(100.0)
		paymentMethod, _ := vo.NewPaymentMethod("credit_card")
		status, _ := vo.NewTransactionStatus("verified")
		expiresAt := time.Now().Add(10 * time.Minute) // Not expired

		cachedTx := &transaction.Transaction{
			ID:            transactionID,
			UserID:        userID,
			Amount:        amount,
			PaymentMethod: paymentMethod,
			Status:        status,
			ExpiresAt:     expiresAt,
		}

		userWallet := &wallet.Wallet{
			ID:      userID,
			Balance: amount,
		}

		// Mock cache get (hit)
		cacheService.On("Get", mock.Anything, "transaction:6", mock.AnythingOfType("*transaction.Transaction")).
			Run(func(args mock.Arguments) {
				tx := args.Get(2).(*transaction.Transaction)
				*tx = *cachedTx
			}).
			Return(nil).Once()

		// Mock logger info
		logger.On("Info", mock.Anything, mock.Anything).Return().Maybe()

		// Mock wallet retrieval
		walletRepo.On("FindById", userID).Return(userWallet, nil).Once()
		// Mock wallet update in transaction
		walletRepo.On("Update", mock.Anything).Return(errors.New("database error")).Once()
		// Mock transaction
		dbTx.On("DoInTransaction", mock.AnythingOfType("func(domain.Repository) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(0).(func(domain.Repository) error)
				err := fn(NewMockRepository(userRepo, walletRepo, transactionRepo))
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
			}).
			Return(errors.New("database error")).Once()
		// Execute
		_, _, err := walletUsecase.ConfirmTopup(transactionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		// Verify mocks
		cacheService.AssertExpectations(t)
		walletRepo.AssertExpectations(t)
		dbTx.AssertExpectations(t)
	})
}
