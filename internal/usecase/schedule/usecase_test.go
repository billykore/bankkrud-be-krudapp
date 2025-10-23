package schedule

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/schedule"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

func TestGetSchedules_Success(t *testing.T) {
	var (
		ctx = context.WithValue(context.Background(), user.ContextKey, user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			FirstName:   "John",
			LastName:    "Doe",
		})
		scheduleRepo = schedule.NewMockRepository(t)
		uc           = NewUsecase(scheduleRepo)
	)

	scheduleRepo.EXPECT().GetAll(ctx, mock.Anything).
		Return([]schedule.Schedule{
			{
				UUID:              "schedule-uuid-123",
				Name:              "Monthly Rent",
				TransactionType:   "transfer",
				SourceAccount:     "123456789",
				DestinationNumber: "987654321",
				Amount:            1500,
				Period:            schedule.MonthlyPeriod,
			},
		}, nil)

	ss, err := uc.GetSchedules(ctx, &GetSchedulesRequest{
		TransactionType: "transfer",
	})

	assert.NoError(t, err)
	assert.NotNil(t, ss)
	assert.Len(t, ss, 1)

	scheduleRepo.AssertExpectations(t)
}

func TestGetScheduleByUUID_Success(t *testing.T) {
	var (
		ctx = context.WithValue(context.Background(), user.ContextKey, user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			FirstName:   "John",
			LastName:    "Doe",
		})
		scheduleRepo = schedule.NewMockRepository(t)
		uc           = NewUsecase(scheduleRepo)
	)

	scheduleRepo.EXPECT().GetByUUID(ctx, "schedule-uuid-123").
		Return(schedule.Schedule{
			UUID:              "schedule-uuid-123",
			Name:              "Monthly Rent",
			TransactionType:   "transfer",
			SourceAccount:     "123456789",
			DestinationNumber: "987654321",
			Amount:            1500,
			Period:            schedule.MonthlyPeriod,
		}, nil)

	s, err := uc.GetScheduleByUUID(ctx, &GetScheduleByUUIDRequest{
		UUID: "schedule-uuid-123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, s.UUID, "schedule-uuid-123")

	scheduleRepo.AssertExpectations(t)
}

func TestCreate_Success(t *testing.T) {
	var (
		ctx = context.WithValue(context.Background(), user.ContextKey, user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			FirstName:   "John",
			LastName:    "Doe",
		})
		scheduleRepo = schedule.NewMockRepository(t)
		uc           = NewUsecase(scheduleRepo)
	)

	scheduleRepo.EXPECT().Create(ctx, mock.Anything).
		Return(nil)

	res, err := uc.CreateSchedule(ctx, &CreateScheduleRequest{
		Name:              "Monthly Rent",
		TransactionType:   "transfer",
		SourceAccount:     "123456789",
		DestinationNumber: "987654321",
		Amount:            1500,
		Period:            schedule.MonthlyPeriod,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.UUID)

	scheduleRepo.AssertExpectations(t)
}

func TestCreate_Failed(t *testing.T) {
	var (
		ctx = context.WithValue(context.Background(), user.ContextKey, user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			FirstName:   "John",
			LastName:    "Doe",
		})
		scheduleRepo = schedule.NewMockRepository(t)
		uc           = NewUsecase(scheduleRepo)
	)

	scheduleRepo.EXPECT().Create(ctx, mock.Anything).
		Return(errors.New("mock error"))

	res, err := uc.CreateSchedule(ctx, &CreateScheduleRequest{
		Name:              "Monthly Rent",
		TransactionType:   "transfer",
		SourceAccount:     "123456789",
		DestinationNumber: "987654321",
		Amount:            1500,
		Period:            schedule.MonthlyPeriod,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, pkgerror.InternalServerError().SetMsg("Failed to create schedule"), err)

	scheduleRepo.AssertExpectations(t)
}
