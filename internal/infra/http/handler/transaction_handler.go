package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/response"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/validation"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/transaction"
)

type TransactionHandler struct {
	va *validation.Validator
	uc *transaction.Usecase
}

func NewTransactionHandler(va *validation.Validator, uc *transaction.Usecase) *TransactionHandler {
	return &TransactionHandler{
		va: va,
		uc: uc,
	}
}

// GetTransactions swaggo annotation.
//
//	@Summary		Get transactions
//	@Description	Get transactions with filters
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string	true	"Authorization token"
//	@Param			uuid				query		string	false	"Transaction UUID"
//	@Param			transaction_type	query		string	false	"Transaction Type"
//	@Param			source_account		query		string	false	"Source Account UUID"
//	@Param			target_account		query		string	false	"Target Account UUID"
//	@Param			status				query		string	false	"Transaction Status"
//	@Success		200					{object}	response.Response
//	@Failure		400					{object}	response.Response
//	@Failure		404					{object}	response.Response
//	@Failure		500					{object}	response.Response
//	@Router			/transactions [get]
func (h *TransactionHandler) GetTransactions(ctx echo.Context) error {
	req := new(transaction.GetTransactionsRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.GetTransactions(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// GetTransaction swaggo annotation.
//
//	@Summary		Get transaction by UUID
//	@Description	Get transaction by UUID
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization token"
//	@Param			uuid			path		string	true	"Transaction UUID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/transactions/{uuid} [get]
func (h *TransactionHandler) GetTransaction(ctx echo.Context) error {
	req := new(transaction.GetTransactionRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.GetTransactionByUUID(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
