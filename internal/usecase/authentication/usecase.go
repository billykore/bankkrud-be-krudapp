package authentication

import (
	"context"
	"errors"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

// Usecase implements the authentication usecase.
type Usecase struct {
	userRepo user.Repository
	authSvc  user.AuthService
}

func NewUsecase(userRepo user.Repository, authSvc user.AuthService) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

func (uc *Usecase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	l := log.WithContext(ctx, "Login")

	usr, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		l.Error().Err(err).
			Str("username", req.Username).
			Msg("User not found")
		return nil, pkgerror.NotFound().SetMsg("User not found")
	}

	err = uc.authSvc.ValidatePassword(req.Password, usr.Password)
	if err != nil {
		l.Error().Err(err).
			Str("username", req.Username).
			Msg("Invalid password")
		return nil, pkgerror.Unauthorized().SetMsg("Invalid password")
	}

	cacheToken, err := uc.userRepo.GetToken(ctx, usr.Username)
	if err != nil && !errors.Is(err, user.ErrTokenNotFound) {
		l.Error().Err(err).
			Str("username", req.Username).
			Msg("Failed to get token")
		return nil, pkgerror.InternalServerError().SetMsg("Failed to get token")
	}
	if cacheToken.Value != "" {
		return &LoginResponse{
			Username:        usr.Username,
			Email:           usr.Email,
			PhoneNumber:     usr.PhoneNumber,
			FullName:        usr.FullName(),
			Token:           cacheToken.Value,
			ExpiredDuration: cacheToken.ExpiredDuration(),
		}, nil
	}

	token, err := uc.authSvc.GenerateToken(usr)
	if err != nil {
		l.Error().Err(err).
			Str("username", req.Username).
			Msg("Failed to generate token")
		return nil, pkgerror.InternalServerError().SetMsg("Failed to generate token")
	}

	err = uc.userRepo.SaveToken(ctx, usr.Username, token)
	if err != nil {
		l.Error().Err(err).
			Str("username", req.Username).
			Msg("Failed to save token")
	}

	return &LoginResponse{
		Username:        usr.Username,
		Email:           usr.Email,
		PhoneNumber:     usr.PhoneNumber,
		FullName:        usr.FullName(),
		Token:           token.Value,
		ExpiredDuration: token.ExpiredDuration(),
	}, nil
}
