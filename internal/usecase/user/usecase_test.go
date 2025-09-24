package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

func TestGetByUsername_Success(t *testing.T) {
	var (
		userRepo = user.NewMockRepository(t)
		uc       = NewUsecase(userRepo)
	)

	userRepo.EXPECT().GetByUsername(mock.Anything, "johndoe").
		Return(user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			FirstName:   "John",
			LastName:    "Doe",
		}, nil)

	res, err := uc.GetByUsername(context.Background(), &GetByUsernameRequest{
		Fields: "username",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "johndoe", res.Username)
	assert.Equal(t, "johndoe@example.com", res.Email)
	assert.Equal(t, "1234567890", res.PhoneNumber)
	assert.Equal(t, "John Doe", res.FullName)
	userRepo.AssertExpectations(t)
}

func TestGetByUsername_RepositoryError(t *testing.T) {
	var (
		userRepo = user.NewMockRepository(t)
		uc       = NewUsecase(userRepo)
	)

	userRepo.EXPECT().GetByUsername(mock.Anything, "johndoe").
		Return(user.User{}, errors.New("mock error"))

	res, err := uc.GetByUsername(context.Background(), &GetByUsernameRequest{
		Fields: "username",
	})

	assert.Equal(t, err, pkgerror.NotFound().SetMsg("User not found"))
	assert.Nil(t, res)

	userRepo.AssertExpectations(t)
}
