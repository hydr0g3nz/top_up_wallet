package model

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents the transactions table
type Transaction struct {
	gorm.Model
	UserID        uint       `gorm:"not null"`
	Amount        float64    `gorm:"type:decimal(18,2);not null;check:amount > 0"`
	PaymentMethod string     `gorm:"size:50;not null;check:payment_method IN ('credit_card')"`
	Status        string     `gorm:"size:20;not null;check:status IN ('verified','completed','failed','expired')"`
	ExpiresAt     *time.Time `gorm:"not null"`
}
