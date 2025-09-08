package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/response"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/validation"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/transfer"
)

type TransferHandler struct {
	va *validation.Validator
	uc *transfer.Usecase
}

func NewTransferHandler(va *validation.Validator, uc *transfer.Usecase) *TransferHandler {
	return &TransferHandler{
		va: va,
		uc: uc,
	}
}

// Initiate swaggo annotation.
//
//	@Summary		Initiate transfer
//	@Description	Initiate transfer transaction
//	@Tags			transfer
//	@Accept			json
//	@Produce		json
//	@Param			InitiateRequest	body		transfer.InitiateRequest	true	"Initiate Transfer Request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/transfer/initiate [post]
func (h *TransferHandler) Initiate(ctx echo.Context) error {
	req := new(transfer.InitiateRequest)
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

// Process swaggo annotation.
//
//	@Summary		Process transfer
//	@Description	Process transfer transaction
//	@Tags			transfer
//	@Accept			json
//	@Produce		json
//	@Param			ProcessRequest	body		transfer.ProcessRequest	true	"Process Transfer Request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/transfer/process [post]
func (h *TransferHandler) Process(ctx echo.Context) error {
	req := new(transfer.ProcessRequest)
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
