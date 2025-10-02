package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

type Usecase struct {
	userRepo    user.Repository
	authSvc     user.AuthService
	accountRepo account.Repository
}

func NewUsecase(userRepo user.Repository, authSvc user.AuthService, accountRepo account.Repository) *Usecase {
	return &Usecase{
		userRepo:    userRepo,
		authSvc:     authSvc,
		accountRepo: accountRepo,
	}
}

func (uc *Usecase) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	l := log.WithContext(ctx, "Create")

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		l.Error().Err(err).Msg("Invalid date of birth")
		return nil, pkgerror.BadRequest().SetMsg("Invalid date of birth")
	}

	hashedPassword, err := uc.authSvc.HashPassword(req.Password)
	if err != nil {
		l.Error().Err(err).Msg("Error hashing password")
		return nil, pkgerror.InternalServerError().SetMsg("Error creating user")
	}

	acc, err := uc.accountRepo.Create(ctx, req.Username)
	if err != nil {
		l.Error().Err(err).Msg("Error creating account")
		return nil, pkgerror.InternalServerError().SetMsg("Error creating user")
	}

	err = uc.userRepo.Create(ctx, user.User{
		CIF:         acc.CIF,
		UUID:        uuid.New().String(),
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		DateOfBirth: dob,
		Status:      user.StatusActive,
		Password:    hashedPassword,
	})
	if err != nil && errors.Is(err, user.ErrDuplicateUserData) {
		l.Error().Err(err).Msg("Duplicate user data")
		return nil, pkgerror.Conflict().SetMsg(err.Error())
	}
	if err != nil {
		l.Error().Err(err).Msg("Error creating user")
		return nil, pkgerror.InternalServerError().SetMsg("Error creating user")
	}

	return &CreateResponse{
		Message: "User registered successfully",
	}, nil
}

func (uc *Usecase) GetByUsername(ctx context.Context, req *GetByUsernameRequest) (*GetByUsernameResponse, error) {
	l := log.WithContext(ctx, "GetByUsername")

	userFromCtx, err := user.FromContext(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Error getting user from context")
		return nil, pkgerror.Unauthorized().SetMsg("User unauthorized")
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
	defaultFields := []string{"username", "first_name", "last_name"}
	if requestFields == "" {
		return defaultFields
	}
	fieldsArray := strings.Split(requestFields, ",")
	return append(defaultFields, fieldsArray...)
}
