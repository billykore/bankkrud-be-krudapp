package repo

import (
	"context"

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
	var saku model.SakuRaya
	res := r.db.WithContext(ctx).
		Where(`"ID" = ?`, id).
		Where(`"STATUS" = ?`, model.SakuStatusOpened).
		First(&saku)
	if err := res.Error; err != nil {
		return pocket.Pocket{}, err
	}
	return pocket.Pocket{
		ID:            saku.ID,
		AccountNumber: saku.CoreCode,
		Name:          saku.Name,
		Status:        saku.Status,
	}, nil
}
