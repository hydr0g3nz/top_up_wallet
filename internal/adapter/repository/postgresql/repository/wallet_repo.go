package repository

import (
	"context"

	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/model"
	IRepository "github.com/hydr0g3nz/wallet_topup_system/internal/domain"
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

func (r *WalletRepository) Update(ctx context.Context, wallet wallet.Wallet) error {
	db := r.getDB(ctx)
	return db.Model(&model.Wallet{}).Where("id = ?", wallet.ID).Updates(wallet.ToNotEmptyValueMap()).Error
}
func (r *WalletRepository) FindById(id uint) (*wallet.Wallet, error) {
	var walletModel model.Wallet
	if err := r.db.Take(&walletModel, id).Error; err != nil {
		return nil, err
	}
	w, err := walletModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return &w, nil
}
func (r *WalletRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := IRepository.GetTx(ctx).(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}
