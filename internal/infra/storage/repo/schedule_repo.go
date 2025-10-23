package repo

import (
	"context"
	"errors"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/schedule"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/storage/model"
	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
	return &ScheduleRepo{
		db: db,
	}
}

func (r *ScheduleRepo) GetAll(ctx context.Context, query schedule.Query) ([]schedule.Schedule, error) {
	m := make([]model.Schedule, 0, query.Limit)
	res := r.db.WithContext(ctx).
		Limit(query.Limit).
		Offset(query.Offset)
	// apply query params
	for k, v := range query.Filter {
		if v != "" {
			res.Where(k, v)
		}
	}
	// finds data
	err := res.Find(&m).Error
	if err != nil {
		return nil, err
	}
	schedules := make([]schedule.Schedule, 0, query.Limit)
	for _, s := range m {
		schedules = append(schedules, schedule.Schedule{
			UUID:              s.UUID,
			Username:          s.Username,
			Name:              s.Name,
			TransactionType:   s.TransactionType,
			SourceAccount:     s.SourceAccount,
			DestinationNumber: s.DestinationNumber,
			Amount:            s.Amount,
			Period:            s.Period,
			Status:            s.Status,
		})
	}
	return schedules, nil
}

func (r *ScheduleRepo) GetByUUID(ctx context.Context, uuid string) (schedule.Schedule, error) {
	var m model.Schedule
	err := r.db.WithContext(ctx).Debug().
		Where("uuid = ?", uuid).
		First(&m).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return schedule.Schedule{}, schedule.ErrNotFound
	}
	if err != nil {
		return schedule.Schedule{}, err
	}
	return schedule.Schedule{
		UUID:              m.UUID,
		Username:          m.Username,
		Name:              m.Name,
		TransactionType:   m.TransactionType,
		SourceAccount:     m.SourceAccount,
		DestinationNumber: m.DestinationNumber,
		Amount:            m.Amount,
		Period:            m.Period,
		Status:            m.Status,
	}, nil
}

func (r *ScheduleRepo) Create(ctx context.Context, schedule schedule.Schedule) error {
	res := r.db.WithContext(ctx).Create(&model.Schedule{
		UUID:              schedule.UUID,
		Username:          schedule.Username,
		Name:              schedule.Name,
		TransactionType:   schedule.TransactionType,
		SourceAccount:     schedule.SourceAccount,
		DestinationNumber: schedule.DestinationNumber,
		Amount:            schedule.Amount,
		Period:            schedule.Period,
		Status:            schedule.Status,
		CronExpression:    "",
	})
	return res.Error
}

func (r *ScheduleRepo) UpdateStatus(ctx context.Context, uuid string, status string) error {
	var m model.Schedule
	res := r.db.WithContext(ctx).Model(&m).
		Where("uuid = ?", uuid).
		Update("status", status)
	return res.Error
}

func (r *ScheduleRepo) Delete(ctx context.Context, uuid string) error {
	var m model.Schedule
	res := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Delete(&m)
	return res.Error
}
