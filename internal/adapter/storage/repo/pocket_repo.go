package repo

import (
	"context"
	"errors"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/storage/model"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/pocket"
	"gorm.io/gorm"
)

type PocketRepo struct {
	db *gorm.DB
}

func NewPocketRepo(db *gorm.DB) *PocketRepo {
	return &PocketRepo{
		db: db,
	}
}

func (r *PocketRepo) GetByAccountNumber(ctx context.Context, accountNumber string) (pocket.Pocket, error) {
	var m model.Pocket
	res := r.db.WithContext(ctx).
		Where("account_number = ?", accountNumber).
		First(&m)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return pocket.Pocket{}, pocket.ErrNotFound
	}
	if res.Error != nil {
		return pocket.Pocket{}, res.Error
	}
	return pocket.Pocket{
		ID:            uint64(m.ID),
		AccountNumber: m.AccountNumber,
		Name:          m.Name,
		Status:        m.Status,
	}, nil
}
