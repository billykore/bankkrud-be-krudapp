package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/response"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/validation"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/user"
)

type UserHandler struct {
	va *validation.Validator
	uc *user.Usecase
}

func NewUserHandler(va *validation.Validator, uc *user.Usecase) *UserHandler {
	return &UserHandler{
		va: va,
		uc: uc,
	}
}

// Create swaggo annotation.
//
//	@Summary		Create a new user
//	@Description	Create a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		user.CreateRequest	true	"Create request"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/users [post]
func (h *UserHandler) Create(ctx echo.Context) error {
	req := new(user.CreateRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	res, err := h.uc.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(res))
}

// GetByUsername swaggo annotation.
//
//	@Summary		Get logged in user
//	@Description	Get logged in user details
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization token"
//	@Param			fields			query		string	false	"Fields"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/users/me [get]
func (h *UserHandler) GetByUsername(ctx echo.Context) error {
	req := new(user.GetByUsernameRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	res, err := h.uc.GetByUsername(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(res))
}
