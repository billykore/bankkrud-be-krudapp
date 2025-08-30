package tapmoney

import (
	"context"
	"errors"

	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/account"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/cbs"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/payment"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/pocket"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/transaction"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/log"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/pkgerror"
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
	log         log.Logger
	cbs         cbs.Service
	txRepo      transaction.Repository
	pocketRepo  pocket.Repository
	paymentSvc  payment.Service
	accountRepo account.Repository
}

func NewUsecase(
	log log.Logger,
	cbs cbs.Service,
	txRepo transaction.Repository,
	pocketRepo pocket.Repository,
	paymentSvc payment.Service,
	accountRepo account.Repository) *Usecase {
	return &Usecase{
		log:         log,
		cbs:         cbs,
		txRepo:      txRepo,
		pocketRepo:  pocketRepo,
		paymentSvc:  paymentSvc,
		accountRepo: accountRepo,
	}
}

func (uc *Usecase) Inquiry(ctx context.Context, req *InquiryRequest) (*InquiryResponse, error) {
	l := uc.log.Usecase("Inquiry")

	cbsStatus, err := uc.cbs.GetStatus(ctx)
	if err != nil {
		l.Errorf("Failed to Get CBS status: %v", err)
		return nil, pkgerror.InternalServerError()
	}
	if cbsStatus.NotReady() {
		l.Errorf("CBS is not ready for transactions")
		return nil, pkgerror.InternalServerError()
	}

	thePocket, err := uc.pocketRepo.GetByAccountNumber(ctx, req.SourceAccount)
	if err != nil && errors.Is(err, pocket.ErrNotFound) {
		l.Errorf("Failed to Get pocket: %v", err)
		return nil, pkgerror.NotFound().SetMsg("Pocket not found")
	}
	if err != nil {
		l.Errorf("Failed to Get pocket: %v", err)
		return nil, pkgerror.InternalServerError()
	}

	result, err := uc.paymentSvc.Inquiry(ctx, tapMoneyChannel, payment.Bill{
		DestinationAccount: req.CardNumber,
		BillerCode:         tapMoneyBillerCode,
		Amount:             req.Amount,
		SourceAccount:      thePocket.AccountNumber,
	})
	if err != nil {
		l.Errorf("Inquiry to payment service failed: %v", err)
		return nil, pkgerror.BadRequest().SetMsg("Invalid card number")
	}

	tx := transaction.Transaction{
		TransactionType:    tapMoneyTransactionType,
		SourceAccount:      thePocket.AccountNumber,
		DestinationAccount: req.CardNumber,
		Status:             transaction.StatusInquirySuccess,
		PaymentID:          result.ID,
		Amount:             req.Amount,
	}

	err = uc.txRepo.Create(ctx, tx)
	if err != nil {
		l.Errorf("Create transaction failed: %v", err)
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
	l := uc.log.Usecase("Payment")

	cbsStatus, err := uc.cbs.GetStatus(ctx)
	if err != nil {
		l.Errorf("Failed to Get CBS status: %v", err)
		return nil, pkgerror.InternalServerError()
	}
	if cbsStatus.NotReady() {
		l.Errorf("CBS is not ready for transactions")
		return nil, pkgerror.InternalServerError()
	}

	tx, err := uc.txRepo.Get(ctx, req.TransactionID)
	if err != nil {
		l.Errorf("Failed to Get transaction: %v", err)
		return nil, pkgerror.NotFound().SetMsg("Transaction was not found")
	}
	if tx.Status != transaction.StatusPending {
		l.Errorf("Transaction is not pending. Status: %s", tx.Status)
		return nil, pkgerror.BadRequest().SetMsg("Transaction is already processed")
	}

	srcAccount, err := uc.accountRepo.Get(ctx, tx.SourceAccount)
	if err != nil {
		l.Errorf("Failed to Get account: %v", err)
		return nil, pkgerror.NotFound().SetMsg("Source account was not found")
	}
	if !srcAccount.CanTransfer(tx.Amount) {
		l.Errorf("Insufficient balance. Balance: %d, Amount: %d", srcAccount.Balance, tx.Amount)
		return nil, pkgerror.BadRequest().SetMsg("Insufficient balance")
	}

	payRes, err := uc.paymentSvc.Payment(ctx, payment.Bill{
		DestinationAccount: tx.DestinationAccount,
		BillerCode:         tapMoneyBillerCode,
		Amount:             tx.Amount,
		SourceAccount:      tx.SourceAccount,
	})
	if err != nil {
		l.Errorf("Payment to payment service failed: %v", err)
		return nil, pkgerror.InternalServerError()
	}

	err = uc.txRepo.Update(ctx, transaction.Transaction{
		UUID:   req.TransactionID,
		Status: transaction.StatusSuccess,
	})
	if err != nil {
		l.Errorf("Update transaction failed: %v", err)
		return nil, pkgerror.InternalServerError()
	}

	return &PaymentResponse{
		TransactionID: tx.UUID,
		Message:       SuccessfulMessage,
		Status:        payRes.Status,
		Amount:        tx.Amount,
		CardNumber:    tx.DestinationAccount,
		Notes:         tx.Notes,
		Fee:           tx.Fee,
	}, nil
}
