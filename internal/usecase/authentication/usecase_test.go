package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
)

func TestLogin_Success(t *testing.T) {
	var (
		userRepo = user.NewMockRepository(t)
		authSvc  = user.NewMockAuthService(t)
		uc       = NewUsecase(userRepo, authSvc)
	)

	userRepo.EXPECT().GetByUsername(mock.Anything, "johndoe").
		Return(user.User{
			Username:    "johndoe",
			Email:       "johndoe@example.com",
			PhoneNumber: "1234567890",
			Password:    "hashed-password",
		}, nil)

	authSvc.EXPECT().ValidatePassword(mock.Anything, mock.Anything).
		Return(nil)

	userRepo.EXPECT().GetToken(mock.Anything, mock.Anything).
		Return(user.Token{}, user.ErrTokenNotFound)

	authSvc.EXPECT().GenerateToken(mock.Anything).
		Return(user.Token{
			Value:     "token-123",
			ExpiresAt: time.Time{},
		}, nil)

	userRepo.EXPECT().SaveToken(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	res, err := uc.Login(context.Background(), &LoginRequest{
		Username: "johndoe",
		Password: "password",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "token-123", res.Token)
	assert.Equal(t, "johndoe", res.Username)
	assert.Equal(t, "johndoe@example.com", res.Email)
	assert.Equal(t, "1234567890", res.PhoneNumber)
}
