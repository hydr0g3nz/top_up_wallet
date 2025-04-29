package repository

import (
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/model"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/wallet"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet wallet.Wallet) error {
	walletModel := model.CreateWalletFromDomain(wallet)
	return r.db.Create(&walletModel).Error
}

func (r *WalletRepository) Update(wallet wallet.Wallet) error {
	return r.db.Updates(wallet.ToNotEmptyValueMap()).Error
}
