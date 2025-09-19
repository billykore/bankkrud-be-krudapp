package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/response"
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
