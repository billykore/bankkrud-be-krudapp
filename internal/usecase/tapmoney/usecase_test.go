package tapmoney

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/log"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/cbs"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/payment"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/pocket"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/transaction"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/pkgerror"
)

func TestInquiry_Success(t *testing.T) {
	var (
		zap        = log.NewZap()
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		pocketRepo = pocket.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		uc         = NewUsecase(zap, cbsService, txRepo, pocketRepo, paymentSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	pocketRepo.EXPECT().Get(mock.Anything, int64(123)).
		Return(pocket.Pocket{
			ID:            123,
			AccountNumber: "001201001479315",
			Name:          "Savings",
			Status:        "Opened",
		}, nil)
	txRepo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(nil)
	paymentSvc.EXPECT().Inquiry(mock.Anything, mock.Anything, mock.Anything).
		Return(payment.Payment{
			ID: "pay-123",
		}, nil)

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber: "6013501000500719",
		PocketID:   123,
		Amount:     10000,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "inq-success", resp.Status)

	t.Log(resp)

	cbsService.AssertExpectations(t)
	pocketRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}

func TestInquiry_GetCbsFailed(t *testing.T) {
	var (
		zap        = log.NewZap()
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		pocketRepo = pocket.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		uc         = NewUsecase(zap, cbsService, txRepo, pocketRepo, paymentSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{}, errors.New("get cbs status failed"))

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber: "6013501000500719",
		PocketID:   123,
		Amount:     10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	pocketRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}

func TestInquiry_GetPocketFailed(t *testing.T) {
	var (
		zap        = log.NewZap()
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		pocketRepo = pocket.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		uc         = NewUsecase(zap, cbsService, txRepo, pocketRepo, paymentSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	pocketRepo.EXPECT().Get(mock.Anything, int64(123)).
		Return(pocket.Pocket{}, errors.New("pocket not found"))

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber: "6013501000500719",
		PocketID:   123,
		Amount:     10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.NotFound().SetMsg("Pocket not found"), err)

	pocketRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}
