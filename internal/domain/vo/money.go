package vo

import (
	"fmt"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain"
)

type Money struct {
	amount float64
}

func NewMoney(amount float64) (Money, error) {
	if amount < 0 {
		return Money{}, domain.ErrNegativeAmount
	}
	return Money{amount: amount}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount}
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.amount < other.amount {
		return Money{}, domain.ErrInsufficientBalance
	}
	return Money{amount: m.amount - other.amount}, nil
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f", m.amount)
}
