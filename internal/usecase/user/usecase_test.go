package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

func TestGetByUsername_Success(t *testing.T) {
	var (
		ctx = context.WithValue(context.Background(), user.ContextKey, user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			FirstName:   "John",
			LastName:    "Doe",
		})
		userRepo    = user.NewMockRepository(t)
		authSvc     = user.NewMockAuthService(t)
		accountRepo = account.NewMockRepository(t)
		uc          = NewUsecase(userRepo, authSvc, accountRepo)
	)

	userRepo.EXPECT().GetFieldsByUsername(mock.Anything, "johndoe", "username", "first_name", "last_name").
		Return(user.User{
			Username:  "johndoe",
			FirstName: "John",
			LastName:  "Doe",
		}, nil)

	res, err := uc.GetByUsername(ctx, &GetByUsernameRequest{
		Fields: "",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "johndoe", res.Username)
	assert.Equal(t, "John Doe", res.FullName)
	userRepo.AssertExpectations(t)
}

func TestGetByUsername_RepositoryError(t *testing.T) {
	var (
		ctx = context.WithValue(context.Background(), user.ContextKey, user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			FirstName:   "John",
			LastName:    "Doe",
		})
		userRepo    = user.NewMockRepository(t)
		authSvc     = user.NewMockAuthService(t)
		accountRepo = account.NewMockRepository(t)
		uc          = NewUsecase(userRepo, authSvc, accountRepo)
	)

	userRepo.EXPECT().GetFieldsByUsername(mock.Anything, "johndoe", "username", "first_name", "last_name").
		Return(user.User{}, errors.New("mock error"))

	res, err := uc.GetByUsername(ctx, &GetByUsernameRequest{
		Fields: "",
	})

	assert.Equal(t, err, pkgerror.NotFound().SetMsg("User not found"))
	assert.Nil(t, res)

	userRepo.AssertExpectations(t)
}
