package user

import (
	"context"
	"strings"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

type Usecase struct {
	userRepo user.Repository
}

func NewUsecase(userRepo user.Repository) *Usecase {
	return &Usecase{
		userRepo: userRepo,
	}
}

func (uc *Usecase) GetByUsername(ctx context.Context, req *GetByUsernameRequest) (*GetByUsernameResponse, error) {
	l := log.WithContext(ctx, "GetByUsername")

	userFromCtx, err := user.FromContext(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Error getting user from context")
		return nil, pkgerror.NotFound().SetMsg("User not found")
	}

	userFromRepo, err := uc.userRepo.GetFieldsByUsername(ctx, userFromCtx.Username,
		parseFields(req.Fields)...)
	if err != nil {
		l.Error().Err(err).
			Str("username", userFromCtx.Username).
			Msg("User not found")
		return nil, pkgerror.NotFound().SetMsg("User not found")
	}

	return &GetByUsernameResponse{
		CIF:         userFromRepo.CIF,
		Username:    userFromRepo.Username,
		FullName:    userFromRepo.FullName(),
		Email:       userFromRepo.Email,
		PhoneNumber: userFromRepo.PhoneNumber,
		Address:     userFromRepo.Address,
		DateOfBirth: userFromRepo.DateOfBirth,
		LastLogin:   userFromRepo.LastLogin,
	}, nil
}

func parseFields(requestFields string) []string {
	if requestFields == "" {
		return []string{}
	}
	return strings.Split(requestFields, ",")
}
