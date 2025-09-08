package tapmoney

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/payment"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

const (
	tapMoneyChannelID       = "01"
	tapMoneyBillerCode      = "99999"
	tapMoneyTransactionType = "tapmoney"
)

// tapMoneyChannel represents the payment channel for Tap Money transactions.
var tapMoneyChannel = payment.Channel{
	ID: tapMoneyChannelID,
}

// Usecase defines the use case for handling TapMoney transactions.
type Usecase struct {
	cbs         cbs.Service
	txRepo      transaction.Repository
	paymentSvc  payment.Service
	accountRepo account.Service
}

func NewUsecase(
	cbs cbs.Service,
	txRepo transaction.Repository,
	paymentSvc payment.Service,
	accountRepo account.Service) *Usecase {
	return &Usecase{
		cbs:         cbs,
		txRepo:      txRepo,
		paymentSvc:  paymentSvc,
		accountRepo: accountRepo,
	}
}

func (uc *Usecase) Inquiry(ctx context.Context, req *InquiryRequest) (*InquiryResponse, error) {
	l := log.With().Ctx(ctx).Str("usecase", "Inquiry").Logger()

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

	srcAccount, err := uc.accountRepo.Get(ctx, req.SourceAccount)
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

	result, err := uc.paymentSvc.Inquiry(ctx, tapMoneyChannel, payment.Bill{
		DestinationAccount: req.CardNumber,
		BillerCode:         tapMoneyBillerCode,
		Amount:             req.Amount,
		SourceAccount:      req.SourceAccount,
	})
	if err != nil {
		l.Error().Err(err).Msg("Inquiry to payment service failed")
		return nil, pkgerror.BadRequest().SetMsg("Invalid card number")
	}

	tx := transaction.Transaction{
		UUID:               uuid.New().String(),
		TransactionType:    tapMoneyTransactionType,
		SourceAccount:      req.SourceAccount,
		DestinationAccount: req.CardNumber,
		Status:             transaction.StatusInquirySuccess,
		PaymentID:          result.ID,
		Amount:             req.Amount,
	}

	err = uc.txRepo.Create(ctx, tx)
	if err != nil {
		l.Error().Err(err).Msg("Create transaction failed")
		return nil, pkgerror.InternalServerError()
	}

	return &InquiryResponse{
		TransactionID: tx.UUID,
		PaymentID:     result.ID,
		Status:        tx.Status,
		CardNumber:    req.CardNumber,
		SourceAccount: req.SourceAccount,
		Amount:        tx.Amount,
	}, nil
}

func (uc *Usecase) Payment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	l := log.With().Ctx(ctx).Str("usecase", "Payment").Logger()

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
		l.Error().Err(err).Msg("Transaction was not found")
		return nil, pkgerror.NotFound().SetMsg("Transaction was not found")
	}
	if tx.Status != transaction.StatusPending {
		l.Error().Err(err).
			Str("transaction_status", tx.Status).
			Msg("Transaction is already processed")
		return nil, pkgerror.BadRequest().SetMsg("Transaction is already processed")
	}

	srcAccount, err := uc.accountRepo.Get(ctx, tx.SourceAccount)
	if err != nil {
		l.Error().Err(err).Msg("Source account was not found")
		return nil, pkgerror.NotFound().SetMsg("Source account was not found")
	}
	if !srcAccount.CanTransfer(tx.Amount) {
		l.Error().
			Int64("account_balance", srcAccount.Balance).
			Int64("request_amount", tx.Amount).
			Msg("Insufficient balance")
		return nil, pkgerror.BadRequest().SetMsg("Insufficient balance")
	}

	payResp, err := uc.paymentSvc.Payment(ctx, payment.Bill{
		DestinationAccount: tx.DestinationAccount,
		BillerCode:         tapMoneyBillerCode,
		Amount:             tx.Amount,
		SourceAccount:      tx.SourceAccount,
	})
	if err != nil {
		l.Error().Err(err).Msg("Payment to payment service failed")
		return nil, pkgerror.InternalServerError()
	}

	tx.Status = transaction.StatusInquirySuccess
	tx.PaymentID = payResp.ID

	err = uc.txRepo.Update(ctx, tx)
	if err != nil {
		l.Error().Err(err).Msg("Update transaction failed")
		return nil, pkgerror.InternalServerError()
	}

	return &PaymentResponse{
		TransactionID: tx.UUID,
		Message:       SuccessfulMessage,
		Status:        tx.Status,
		Amount:        tx.Amount,
		CardNumber:    tx.DestinationAccount,
		Notes:         tx.Notes,
		Fee:           tx.Fee,
	}, nil
}
