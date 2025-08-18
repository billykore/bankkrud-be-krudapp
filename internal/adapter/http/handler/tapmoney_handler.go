package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/http/response"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/validation"
	"go.bankkrud.com/backend/svc/tapmoney/internal/usecase/tapmoney"
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

// Inquiry swaggo annotation.
//
//	@Summary		Inquiry TapMoney transaction
//	@Description	Inquiry TapMoney transaction
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Param			InquiryRequest	body		tapmoney.InquiryRequest	true	"Inquiry Request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/tapmoney/inquiry [post]
func (h *TapMoneyHandler) Inquiry(ctx echo.Context) error {
	req := new(tapmoney.InquiryRequest)
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.Inquiry(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
