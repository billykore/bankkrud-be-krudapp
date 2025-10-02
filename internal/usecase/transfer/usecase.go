package transfer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transfer"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

const (
	transferTransactionType = "transfer"
)

// Usecase defines the use case for handling transfers.
type Usecase struct {
	cbsSvc      cbs.Service
	txRepo      transaction.Repository
	accountRepo account.Repository
	transferSvc transfer.Service
}

func NewUsecase(
	cbsSvc cbs.Service,
	txRepo transaction.Repository,
	accountRepo account.Repository,
	transferSvc transfer.Service,
) *Usecase {
	return &Usecase{
		cbsSvc:      cbsSvc,
		txRepo:      txRepo,
		accountRepo: accountRepo,
		transferSvc: transferSvc,
	}
}

func (uc *Usecase) Initiate(ctx context.Context, req *InitiateRequest) (*InitiateResponse, error) {
	l := log.WithContext(ctx, "Initiate")

	cbsStatus, err := uc.cbsSvc.GetStatus(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Failed to Get CBS status")
		return nil, pkgerror.InternalServerError()
	}
	if cbsStatus.NotReady() {
		l.Error().
			Bool("is_eod", cbsStatus.IsEOD).
			Bool("is_stand_in", cbsStatus.IsStandIn).
			Msg("CBS is not ready for transactions")
		return nil, pkgerror.InternalServerError()
	}

	srcAccount, err := uc.accountRepo.Get(ctx, req.SourceAccount)
	if err != nil {
		l.Error().Err(err).
			Str("account_number", req.SourceAccount).
			Msg("Failed to get account")
		return nil, pkgerror.InternalServerError()
	}
	if !srcAccount.CanTransfer(req.Amount) {
		l.Error().
			Int64("account_balance", srcAccount.Balance).
			Int64("request_amount", req.Amount).
			Msg("Insufficient balance")
		return nil, pkgerror.BadRequest().SetMsg("Insufficient balance")
	}

	destAccount, err := uc.accountRepo.Get(ctx, req.DestinationAccount)
	if err != nil {
		l.Error().Err(err).
			Str("account_number", req.DestinationAccount).
			Msg("Failed to get account")
		return nil, pkgerror.InternalServerError()
	}

	userFromCtx, err := user.FromContext(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Error getting user from context")
		return nil, pkgerror.Unauthorized().SetMsg("User unauthorized")
	}

	tx := transaction.Transaction{
		UUID:               uuid.New().String(),
		SourceAccount:      srcAccount.AccountNumber,
		DestinationAccount: destAccount.AccountNumber,
		TransactionType:    transferTransactionType,
		Status:             transaction.StatusInitiated,
		Amount:             req.Amount,
		Username:           userFromCtx.Username,
		Note:               req.Note,
	}

	err = uc.txRepo.Create(ctx, tx)
	if err != nil {
		l.Error().Err(err).Msg("Failed to create transaction")
		return nil, pkgerror.InternalServerError()
	}

	return &InitiateResponse{
		UUID:   tx.UUID,
		Status: tx.Status,
	}, nil
}

func (uc *Usecase) Process(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	l := log.WithContext(ctx, "Process")

	cbsStatus, err := uc.cbsSvc.GetStatus(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Failed to Get CBS status")
		return nil, pkgerror.InternalServerError()
	}
	if cbsStatus.NotReady() {
		l.Error().
			Bool("is_eod", cbsStatus.IsEOD).
			Bool("is_stand_in", cbsStatus.IsStandIn).
			Msg("CBS is not ready for transactions")
		return nil, pkgerror.InternalServerError()
	}

	tx, err := uc.txRepo.GetByUUID(ctx, req.UUID)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get transaction")
		return nil, pkgerror.InternalServerError()
	}
	if tx.Status != transaction.StatusInitiated {
		l.Error().
			Str("uuid", req.UUID).
			Str("status", tx.Status).
			Msg("Transaction is not in a valid state to be processed")
		return nil, pkgerror.Conflict().SetMsg("Transaction is not in a valid state to be processed")
	}

	res, err := uc.transferSvc.Transfer(
		ctx,
		tx.SourceAccount,
		tx.DestinationAccount,
		req.Amount,
		makeTransferRemark(tx.SourceAccount, tx.DestinationAccount, tx.UUID),
	)
	if err != nil {
		l.Error().Err(err).Msg("Failed to transfer amount")
		return nil, pkgerror.InternalServerError()
	}

	// Update transaction status to success
	tx.Status = transaction.StatusCompleted
	tx.TransactionReference = res.TransactionReference

	err = uc.txRepo.Update(ctx, tx)
	if err != nil {
		l.Error().Err(err).
			Str("transaction_id", req.UUID).
			Msg("Failed to update transaction status")
		return nil, pkgerror.InternalServerError()
	}

	return &ProcessResponse{
		UUID:   tx.UUID,
		Status: tx.Status,
	}, nil
}

// makeTransferRemark creates a remark for the transfer transaction.
func makeTransferRemark(srcAccount, destAccount, uuid string) string {
	return fmt.Sprintf("TRF %s %s BNKKRD %s", srcAccount, destAccount, uuid)
}
