package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/response"
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

// GetByUsername swaggo annotation.
//
//	@Summary		Get user by username
//	@Description	Get user by username
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string	true	"Username"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Failure		500			{object}	response.Response
//	@Router			/users/{username} [get]
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
