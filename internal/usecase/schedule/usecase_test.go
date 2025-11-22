package schedule

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/schedule"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
)

func TestGetSchedules(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *GetSchedulesRequest
	}
	tests := []struct {
		name    string
		args    args
		mock    func(repo *schedule.MockRepository)
		want    *GetSchedulesResponse
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx: context.WithValue(context.Background(), user.ContextKey, user.User{
					Username:    "johndoe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
					FirstName:   "John",
					LastName:    "Doe",
				}),
				req: &GetSchedulesRequest{
					TransactionType: "transfer",
					Limit:           1,
				},
			},
			mock: func(repo *schedule.MockRepository) {
				repo.EXPECT().GetAll(mock.Anything, mock.Anything).
					Return([]schedule.Schedule{
						{
							ID:                1,
							UUID:              "schedule-uuid-123",
							Name:              "Monthly Rent",
							TransactionType:   "transfer",
							SourceAccount:     "123456789",
							DestinationNumber: "987654321",
							Amount:            1500,
							Period:            schedule.MonthlyPeriod,
						},
						{
							ID:                2,
							UUID:              "schedule-uuid-456",
							Name:              "Weekly Groceries",
							TransactionType:   "transfer",
							SourceAccount:     "123456789",
							DestinationNumber: "987654321",
							Amount:            500,
							Period:            schedule.WeeklyPeriod,
						},
					}, nil)
			},
			want: &GetSchedulesResponse{
				Schedules: []*GetScheduleResponse{
					{
						UUID:              "schedule-uuid-123",
						Name:              "Monthly Rent",
						TransactionType:   "transfer",
						SourceAccount:     "123456789",
						DestinationNumber: "987654321",
						Amount:            1500,
						Period:            schedule.MonthlyPeriod,
					},
				},
				StartID: 1,
				NextID:  2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			scheduleRepo := schedule.NewMockRepository(t)
			uc := NewUsecase(scheduleRepo)
			if tt.mock != nil {
				tt.mock(scheduleRepo)
			}
			got, err := uc.GetSchedules(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetSchedules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetScheduleByUUID(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *GetScheduleByUUIDRequest
	}
	tests := []struct {
		name    string
		args    args
		mock    func(repo *schedule.MockRepository)
		want    *GetScheduleResponse
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx: context.WithValue(context.Background(), user.ContextKey, user.User{
					Username:    "johndoe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
					FirstName:   "John",
					LastName:    "Doe",
				}),
				req: &GetScheduleByUUIDRequest{
					UUID: "schedule-uuid-123",
				},
			},
			mock: func(repo *schedule.MockRepository) {
				repo.EXPECT().GetByUUID(mock.Anything, "schedule-uuid-123").
					Return(schedule.Schedule{
						UUID:              "schedule-uuid-123",
						Name:              "Monthly Rent",
						TransactionType:   "transfer",
						SourceAccount:     "123456789",
						DestinationNumber: "987654321",
						Amount:            1500,
						Period:            schedule.MonthlyPeriod,
					}, nil)
			},
			want: &GetScheduleResponse{
				UUID:              "schedule-uuid-123",
				Name:              "Monthly Rent",
				TransactionType:   "transfer",
				SourceAccount:     "123456789",
				DestinationNumber: "987654321",
				Amount:            1500,
				Period:            schedule.MonthlyPeriod,
			},
			wantErr: false,
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.WithValue(context.Background(), user.ContextKey, user.User{
					Username:    "johndoe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
					FirstName:   "John",
					LastName:    "Doe",
				}),
				req: &GetScheduleByUUIDRequest{
					UUID: "schedule-uuid-not-found",
				},
			},
			mock: func(repo *schedule.MockRepository) {
				repo.EXPECT().GetByUUID(mock.Anything, "schedule-uuid-not-found").
					Return(schedule.Schedule{}, schedule.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			scheduleRepo := schedule.NewMockRepository(t)
			uc := NewUsecase(scheduleRepo)
			if tt.mock != nil {
				tt.mock(scheduleRepo)
			}
			got, err := uc.GetScheduleByUUID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetScheduleByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateSchedule(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *CreateScheduleRequest
	}
	tests := []struct {
		name    string
		args    args
		mock    func(repo *schedule.MockRepository)
		want    *CreateScheduleResponse
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx: context.WithValue(context.Background(), user.ContextKey, user.User{
					Username:    "johndoe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
					FirstName:   "John",
					LastName:    "Doe",
				}),
				req: &CreateScheduleRequest{
					Name:              "Monthly Rent",
					TransactionType:   "transfer",
					SourceAccount:     "123456789",
					DestinationNumber: "987654321",
					Amount:            1500,
					Period:            schedule.MonthlyPeriod,
				},
			},
			mock: func(repo *schedule.MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.Anything).
					Return(nil)
			},
			want: &CreateScheduleResponse{
				UUID: "mock-uuid", // This will be checked loosely as "NotEmpty"
			},
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				ctx: context.WithValue(context.Background(), user.ContextKey, user.User{
					Username:    "johndoe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
					FirstName:   "John",
					LastName:    "Doe",
				}),
				req: &CreateScheduleRequest{
					Name:              "Monthly Rent",
					TransactionType:   "transfer",
					SourceAccount:     "123456789",
					DestinationNumber: "987654321",
					Amount:            1500,
					Period:            schedule.MonthlyPeriod,
				},
			},
			mock: func(repo *schedule.MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.Anything).
					Return(errors.New("mock error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			scheduleRepo := schedule.NewMockRepository(t)
			uc := NewUsecase(scheduleRepo)
			if tt.mock != nil {
				tt.mock(scheduleRepo)
			}
			got, err := uc.CreateSchedule(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.CreateSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.NotEmpty(t, got.UUID)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
