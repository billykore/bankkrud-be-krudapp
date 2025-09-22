package user

import (
	"context"

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

	u, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		l.Error().Err(err).
			Str("username", req.Username).
			Msg("User not found")
		return nil, pkgerror.NotFound().SetMsg("User not found")
	}
	return &GetByUsernameResponse{
		CIF:         u.CIF,
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		FullName:    u.FullName(),
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address,
	}, nil
}
