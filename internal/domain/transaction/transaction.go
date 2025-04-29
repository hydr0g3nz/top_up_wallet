package transaction

import (
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
)

// Transaction represents the transactions table
type Transaction struct {
	ID            uint
	UserID        uint
	Amount        float64
	PaymentMethod vo.PaymentMethod
	Status        vo.TransactionStatus
	ExpiresAt     time.Time
}

func (t Transaction) ToNotEmptyValueMap() map[string]interface{} {
	result := make(map[string]interface{})
	if t.Amount != 0 {
		result["amount"] = t.Amount
	}
	if t.PaymentMethod != "" {
		result["payment_method"] = t.PaymentMethod.String()
	}
	if t.Status != "" {
		result["status"] = t.Status.String()
	}
	if !t.ExpiresAt.IsZero() {
		result["expires_at"] = t.ExpiresAt
	}
	return result
}

type TransactionFilter struct {
	PaymentMethod *vo.PaymentMethod
	Status        *vo.TransactionStatus
	Amount        *vo.Money
	ExpiredAt     *time.Time
}
