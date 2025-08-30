package repo

import (
	"context"
	"errors"

	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/storage/model"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/pocket"
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

func (r *PocketRepo) Get(ctx context.Context, id int64) (pocket.Pocket, error) {
	var m model.Pocket
	res := r.db.WithContext(ctx).
		Where("id = ?", id).
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
