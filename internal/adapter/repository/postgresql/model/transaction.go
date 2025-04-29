package model

import (
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/transaction"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
	"gorm.io/gorm"
)

// Transaction represents the transactions table
type Transaction struct {
	gorm.Model
	UserID        uint      `gorm:"not null"`
	Amount        float64   `gorm:"type:decimal(18,2);not null;check:amount > 0"`
	PaymentMethod string    `gorm:"size:50;not null;check:payment_method IN ('credit_card')"`
	Status        string    `gorm:"size:20;not null;check:status IN ('verified','completed','failed','expired')"`
	ExpiresAt     time.Time `gorm:"not null"`
}

func (t Transaction) ToDomain() transaction.Transaction {
	return transaction.Transaction{
		ID:            t.ID,
		UserID:        t.UserID,
		Amount:        t.Amount,
		PaymentMethod: vo.PaymentMethod(t.PaymentMethod),
		Status:        vo.TransactionStatus(t.Status),
		ExpiresAt:     t.ExpiresAt,
	}
}
func CreateTransactionFromDomain(t transaction.Transaction) Transaction {
	return Transaction{
		Model:         gorm.Model{ID: t.ID},
		UserID:        t.UserID,
		Amount:        t.Amount,
		PaymentMethod: t.PaymentMethod.String(),
		Status:        t.Status.String(),
		ExpiresAt:     t.ExpiresAt,
	}
}
