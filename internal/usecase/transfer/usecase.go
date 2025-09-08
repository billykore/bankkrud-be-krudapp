package transfer

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

const (
	transferTransactionType = "transfer"
)

type Usecase struct {
	cbs        cbs.Service
	txRepo     transaction.Repository
	accountSvc account.Service
}

func NewUsecase(
	cbs cbs.Service,
	txRepo transaction.Repository,
	accountSvc account.Service,
) *Usecase {
	return &Usecase{
		cbs:        cbs,
		txRepo:     txRepo,
		accountSvc: accountSvc,
	}
}

func (uc *Usecase) Initiate(ctx context.Context, req *InitiateRequest) (*InitiateResponse, error) {
	l := log.With().Ctx(ctx).Str("usecase", "Initiate").Logger()

	cbsStatus, err := uc.cbs.GetStatus(ctx)
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

	srcAccount, err := uc.accountSvc.Get(ctx, req.SourceAccount)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get pocket")
		return nil, pkgerror.InternalServerError()
	}
	if !srcAccount.CanTransfer(req.Amount) {
		l.Error().
			Int64("account_balance", srcAccount.Balance).
			Int64("request_amount", req.Amount).
			Msg("Insufficient balance")
		return nil, pkgerror.BadRequest().SetMsg("Insufficient balance")
	}

	destAccount, err := uc.accountSvc.Get(ctx, req.DestinationAccount)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get pocket")
		return nil, pkgerror.InternalServerError()
	}

	tx := transaction.Transaction{
		UUID:               uuid.New().String(),
		SourceAccount:      srcAccount.AccountNumber,
		DestinationAccount: destAccount.AccountNumber,
		TransactionType:    transferTransactionType,
		Status:             transaction.StatusInquirySuccess,
		Amount:             req.Amount,
	}

	err = uc.txRepo.Create(ctx, tx)
	if err != nil {
		l.Error().Err(err).Msg("Failed to create transaction")
		return nil, pkgerror.InternalServerError()
	}

	return &InitiateResponse{
		TransactionID: tx.UUID,
		Status:        tx.Status,
	}, nil
}

func (uc *Usecase) Process(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	l := log.With().Ctx(ctx).Str("usecase", "Process").Logger()

	cbsStatus, err := uc.cbs.GetStatus(ctx)
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

	tx, err := uc.txRepo.Get(ctx, req.TransactionID)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get transaction")
		return nil, pkgerror.InternalServerError()
	}
	if tx.Status != transaction.StatusInquirySuccess {
		l.Error().
			Str("transaction_id", req.TransactionID).
			Str("transaction_status", tx.Status).
			Msg("Transaction is not in a valid state to be processed")
		return nil, pkgerror.BadRequest().SetMsg("Transaction is not in a valid state to be processed")
	}

	err = uc.accountSvc.Transfer(ctx, tx.SourceAccount, tx.DestinationAccount, req.Amount)
	if err != nil {
		l.Error().Err(err).Msg("Failed to transfer amount")
		return nil, pkgerror.InternalServerError()
	}

	// Update transaction status to success
	tx.Status = transaction.StatusSuccess

	err = uc.txRepo.Update(ctx, tx)
	if err != nil {
		l.Error().Err(err).Msg("Failed to update transaction status")
		return nil, pkgerror.InternalServerError()
	}

	return &ProcessResponse{
		TransactionID: tx.UUID,
		Status:        tx.Status,
	}, nil
}
