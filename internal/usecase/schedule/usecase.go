package schedule

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/schedule"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

type Usecase struct {
	scheduleRepo schedule.Repository
}

func NewUsecase(scheduleRepo schedule.Repository) *Usecase {
	return &Usecase{
		scheduleRepo: scheduleRepo,
	}
}

func (uc *Usecase) GetSchedules(ctx context.Context, req *GetSchedulesRequest) ([]*GetScheduleResponse, error) {
	l := log.WithContext(ctx, "GetSchedules")

	userFromCtx, err := user.FromContext(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Error getting user from context")
		return nil, pkgerror.Unauthorized().SetMsg("User not authorized")
	}

	schedules, err := uc.scheduleRepo.GetAll(ctx, schedule.Query{
		Limit:  req.Limit,
		Offset: req.Offset,
		Filter: schedule.Filter{
			"username":         userFromCtx.Username,
			"transaction_type": req.TransactionType,
			"status":           req.Status,
		},
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed to get schedules")
		return nil, pkgerror.InternalServerError().SetMsg("Failed to get schedules")
	}

	resp := make([]*GetScheduleResponse, 0, req.Limit)
	for _, s := range schedules {
		resp = append(resp, &GetScheduleResponse{
			UUID:              s.UUID,
			Name:              s.Name,
			TransactionType:   s.TransactionType,
			SourceAccount:     s.SourceAccount,
			DestinationNumber: s.DestinationNumber,
			Amount:            s.Amount,
			Period:            s.Period,
			Day:               s.Day,
			Date:              s.Date,
			Status:            s.Status,
			Username:          s.Username,
		})
	}

	return resp, nil
}

func (uc *Usecase) GetScheduleByUUID(ctx context.Context, req *GetScheduleByUUIDRequest) (*GetScheduleResponse, error) {
	l := log.WithContext(ctx, "GetScheduleByUUID")

	s, err := uc.scheduleRepo.GetByUUID(ctx, req.UUID)
	if err != nil && errors.Is(err, schedule.ErrNotFound) {
		l.Error().Err(err).Msg("Failed to get schedule by UUID")
		return nil, pkgerror.NotFound().SetMsg("Schedule not found")
	}
	if err != nil {
		l.Error().Err(err).Msg("Failed to get schedule by UUID")
		return nil, pkgerror.InternalServerError().SetMsg("Failed to get schedule")
	}

	return &GetScheduleResponse{
		UUID:              s.UUID,
		Name:              s.Name,
		TransactionType:   s.TransactionType,
		SourceAccount:     s.SourceAccount,
		DestinationNumber: s.DestinationNumber,
		Amount:            s.Amount,
		Period:            s.Period,
		Day:               s.Day,
		Date:              s.Date,
		Status:            s.Status,
		Username:          s.Username,
	}, nil
}

func (uc *Usecase) CreateSchedule(ctx context.Context, req *CreateScheduleRequest) (*CreateScheduleResponse, error) {
	l := log.WithContext(ctx, "CreateSchedule")

	userFromCtx, err := user.FromContext(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Error getting user from context")
		return nil, pkgerror.Unauthorized().SetMsg("User not authorized")
	}

	s := schedule.Schedule{
		UUID:              uuid.New().String(),
		Username:          userFromCtx.Username,
		Name:              req.Name,
		TransactionType:   req.TransactionType,
		SourceAccount:     req.SourceAccount,
		DestinationNumber: req.DestinationNumber,
		Amount:            req.Amount,
		Period:            req.Period,
		Status:            schedule.StatusActive,
	}

	err = uc.scheduleRepo.Create(ctx, s)
	if err != nil {
		l.Error().Err(err).Msg("Failed to create schedule")
		return nil, pkgerror.InternalServerError().SetMsg("Failed to create schedule")
	}

	return &CreateScheduleResponse{
		UUID: s.UUID,
	}, nil
}

func (uc *Usecase) UpdateScheduleStatus(ctx context.Context, req *UpdateScheduleStatusRequest) (*UpdateScheduleStatusResponse, error) {
	l := log.WithContext(ctx, "UpdateScheduleStatus")

	err := uc.scheduleRepo.UpdateStatus(ctx, req.UUID, req.Status)
	if err != nil {
		l.Error().Err(err).Msg("Failed to update schedule status")
		return nil, pkgerror.InternalServerError().SetMsg("Failed to update schedule status")
	}

	return &UpdateScheduleStatusResponse{
		Message: fmt.Sprintf("Success update schedule status to %s", req.Status),
	}, nil
}

func (uc *Usecase) DeleteSchedule(ctx context.Context, req *DeleteScheduleRequest) (*DeleteScheduleResponse, error) {
	l := log.WithContext(ctx, "DeleteSchedule")

	err := uc.scheduleRepo.Delete(ctx, req.UUID)
	if err != nil {
		l.Error().Err(err).Msg("Failed to delete schedule")
		return nil, pkgerror.InternalServerError().SetMsg("Failed to delete schedule")
	}

	return &DeleteScheduleResponse{
		Message: "Success delete schedule",
	}, nil
}
