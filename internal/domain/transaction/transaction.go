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
