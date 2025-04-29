package model

import "gorm.io/gorm"

// User represents the users table
type User struct {
	gorm.Model
	FirstName string `gorm:"size:50;not null"`
	LastName  string `gorm:"size:50;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Password  string `gorm:"size:255;not null"`
	Phone     string `gorm:"size:20;not null"`
	// relationship
	Transactions []Transaction `gorm:"constraint:OnDelete:CASCADE"` // One-to-Many
	Wallet       Wallet        `gorm:"constraint:OnDelete:CASCADE"` // One-to-One
}
