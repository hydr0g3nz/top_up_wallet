package repository

import (
	"context"

	IRepository "github.com/hydr0g3nz/wallet_topup_system/internal/domain"
	"gorm.io/gorm"
)

// TxManager impl สำหรับ GORM
type txManagerGorm struct {
	db *gorm.DB
}

func NewTxManagerGorm(db *gorm.DB) IRepository.TxManager {
	return &txManagerGorm{db: db}
}

func (tm *txManagerGorm) BeginTx(ctx context.Context) (context.Context, error) {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return IRepository.WithTx(ctx, tx), nil
}

func (tm *txManagerGorm) CommitTx(ctx context.Context) error {
	if tx, ok := IRepository.GetTx(ctx).(*gorm.DB); ok {
		return tx.Commit().Error
	}
	return nil // หรือ error ถ้าไม่มี tx
}

func (tm *txManagerGorm) RollbackTx(ctx context.Context) error {
	if tx, ok := IRepository.GetTx(ctx).(*gorm.DB); ok {
		return tx.Rollback().Error
	}
	return nil
}
