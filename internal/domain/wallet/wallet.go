package wallet

import "github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"

// Wallet represents the wallets table (1-to-1 with User)
type Wallet struct {
	ID      uint
	Balance vo.Money
}
