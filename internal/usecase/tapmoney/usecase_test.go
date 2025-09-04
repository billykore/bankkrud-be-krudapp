package tapmoney

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/payment"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/pocket"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

func TestInquiry_Success(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		pocketRepo  = pocket.NewMockRepository(t)
		paymentSvc  = payment.NewMockService(t)
		accountRepo = account.NewMockRepository(t)
		uc          = NewUsecase(cbsService, txRepo, pocketRepo, paymentSvc, accountRepo)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	pocketRepo.EXPECT().GetByAccountNumber(mock.Anything, "123").
		Return(pocket.Pocket{
			ID:            123,
			AccountNumber: "001201001479315",
			Name:          "Savings",
			Status:        "active",
		}, nil)
	paymentSvc.EXPECT().Inquiry(mock.Anything, mock.Anything, mock.Anything).
		Return(payment.Payment{
			ID: "pay-123",
		}, nil)
	txRepo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(nil)

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "123",
		Amount:        10000,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "inq-success", resp.Status)

	t.Log(resp)

	cbsService.AssertExpectations(t)
	pocketRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestInquiry_GetCbsFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		pocketRepo  = pocket.NewMockRepository(t)
		paymentSvc  = payment.NewMockService(t)
		accountRepo = account.NewMockRepository(t)
		uc          = NewUsecase(cbsService, txRepo, pocketRepo, paymentSvc, accountRepo)
	)

	log.Configure("development")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{}, errors.New("get cbs status failed"))

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "123",
		Amount:        10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	pocketRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}

func TestInquiry_PocketNotFound(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		pocketRepo  = pocket.NewMockRepository(t)
		paymentSvc  = payment.NewMockService(t)
		accountRepo = account.NewMockRepository(t)
		uc          = NewUsecase(cbsService, txRepo, pocketRepo, paymentSvc, accountRepo)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	pocketRepo.EXPECT().GetByAccountNumber(mock.Anything, "321").
		Return(pocket.Pocket{}, pocket.ErrNotFound)

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "321",
		Amount:        10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.NotFound().SetMsg("Pocket not found"), err)

	cbsService.AssertExpectations(t)
	pocketRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestPayment_Success(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		pocketRepo  = pocket.NewMockRepository(t)
		paymentSvc  = payment.NewMockService(t)
		accountRepo = account.NewMockRepository(t)
		uc          = NewUsecase(cbsService, txRepo, pocketRepo, paymentSvc, accountRepo)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	txRepo.EXPECT().Get(mock.Anything, mock.Anything).
		Return(transaction.Transaction{
			UUID:               "trx-123",
			SourceAccount:      "001201001479315",
			DestinationAccount: "6013501000500719",
			Amount:             10000,
			Status:             "pending",
			Notes:              "test",
			Fee:                1500,
		}, nil)
	accountRepo.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       1000000,
			AccountNumber: "001201001479315",
		}, nil)
	paymentSvc.EXPECT().Payment(mock.Anything, mock.Anything).
		Return(payment.Payment{
			Status: "success",
		}, nil)
	txRepo.EXPECT().Update(mock.Anything, mock.Anything).
		Return(nil)

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "6013501000500719", resp.CardNumber)
	assert.Equal(t, "trx-123", resp.TransactionID)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Payment successful", resp.Message)
	assert.Equal(t, int64(10000), resp.Amount)
	assert.Equal(t, int64(1500), resp.Fee)
	assert.Equal(t, "test", resp.Notes)

	t.Log(resp)

	cbsService.AssertExpectations(t)
	pocketRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}
