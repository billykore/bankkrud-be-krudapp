package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/response"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/validation"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/tapmoney"
)

type TapMoneyHandler struct {
	va *validation.Validator
	uc *tapmoney.Usecase
}

func NewTapMoneyHandler(va *validation.Validator, uc *tapmoney.Usecase) *TapMoneyHandler {
	return &TapMoneyHandler{
		va: va,
		uc: uc,
	}
}

// Initiate swaggo annotation.
//
//	@Summary		Initiate TapMoney transaction
//	@Description	Initiate TapMoney transaction
//	@Tags			tapmoney
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Authorization token"
//	@Param			InitiateRequest	body		tapmoney.InitiateRequest	true	"Initiate request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/tapmoney/init [post]
func (h *TapMoneyHandler) Initiate(ctx echo.Context) error {
	req := new(tapmoney.InitiateRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.Initiate(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// Payment swaggo annotation.
//
//	@Summary		Process TapMoney transaction
//	@Description	Process TapMoney transaction
//	@Tags			tapmoney
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Authorization token"
//	@Param			ProcessRequest	body		tapmoney.ProcessRequest	true	"Process request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/tapmoney/{uuid}/process [post]
func (h *TapMoneyHandler) Process(ctx echo.Context) error {
	req := new(tapmoney.ProcessRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.Process(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
