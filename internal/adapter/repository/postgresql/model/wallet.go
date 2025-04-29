package model

import "gorm.io/gorm"

// Wallet represents the wallets table (1-to-1 with User)
type Wallet struct {
	gorm.Model
	Balance float64 `gorm:"type:decimal(18,2);not null;default:0.00"`
}
