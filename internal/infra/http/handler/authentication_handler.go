package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/response"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/validation"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/authentication"
)

type AuthenticationHandler struct {
	va *validation.Validator
	uc *authentication.Usecase
}

func NewAuthenticationHandler(va *validation.Validator, uc *authentication.Usecase) *AuthenticationHandler {
	return &AuthenticationHandler{
		va: va,
		uc: uc,
	}
}

// Login swaggo annotation.
//
//	@Summary		User login
//	@Description	User login
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			LoginRequest	body		authentication.LoginRequest	true	"Login Request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/auth/login [post]
func (h *AuthenticationHandler) Login(ctx echo.Context) error {
	req := new(authentication.LoginRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.Login(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
