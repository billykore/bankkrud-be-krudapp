package tapmoney

import (
	"context"

	"github.com/google/uuid"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/cbs"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/payment"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/pocket"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/transaction"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/log"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/pkgerror"
)

const (
	tapMoneyChannelID  = "01"
	tapMoneyBillerCode = "99999"
)

// tapMoneyChannel represents the payment channel for Tap Money transactions.
var tapMoneyChannel = payment.Channel{
	ID: tapMoneyChannelID,
}

// Usecase defines the use case for handling TapMoney transactions.
type Usecase struct {
	log        log.Logger
	cbs        cbs.Service
	txRepo     transaction.Repository
	pocketRepo pocket.Repository
	paymentSvc payment.Service
}

func NewUsecase(log log.Logger, cbs cbs.Service, txRepo transaction.Repository, pocketRepo pocket.Repository, paymentSvc payment.Service) *Usecase {
	return &Usecase{
		log:        log,
		cbs:        cbs,
		txRepo:     txRepo,
		pocketRepo: pocketRepo,
		paymentSvc: paymentSvc,
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

	thePocket, err := uc.pocketRepo.Get(ctx, req.PocketID)
	if err != nil {
		l.Errorf("Failed to Get pocket: %v", err)
		return nil, pkgerror.NotFound().SetMsg("Pocket not found")
	}

	paymentRes, err := uc.paymentSvc.Inquiry(ctx, tapMoneyChannel, payment.Bill{
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
		UUID:               uuid.NewString(),
		SourceAccount:      thePocket.AccountNumber,
		DestinationAccount: req.CardNumber,
		Amount:             req.Amount,
		Status:             transaction.InquirySuccess,
		PaymentID:          paymentRes.ID,
	}

	err = uc.txRepo.Create(ctx, tx)
	if err != nil {
		l.Errorf("Create transaction failed: %v", err)
		return nil, pkgerror.InternalServerError()
	}

	return &InquiryResponse{
		TransactionID: tx.UUID,
		PaymentID:     paymentRes.ID,
		Status:        tx.Status,
		Amount:        tx.Amount,
		CardNumber:    req.CardNumber,
		PocketID:      req.PocketID,
	}, nil
}
