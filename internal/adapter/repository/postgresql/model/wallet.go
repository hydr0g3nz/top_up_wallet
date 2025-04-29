package model

import (
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/wallet"
	"gorm.io/gorm"
)

// Wallet represents the wallets table (1-to-1 with User)
type Wallet struct {
	gorm.Model
	Balance float64 `gorm:"type:decimal(18,2);not null;default:0.00"`
}

func CreateWalletFromDomain(w wallet.Wallet) Wallet {
	return Wallet{
		Model:   gorm.Model{ID: w.ID},
		Balance: w.Balance.Amount(),
	}
}

func (w Wallet) ToDomain() (wallet.Wallet, error) {
	money, err := vo.NewMoney(w.Balance)
	if err != nil {
		return wallet.Wallet{}, err
	}
	return wallet.Wallet{
		ID:      w.ID,
		Balance: money,
	}, nil
}
