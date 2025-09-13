package transfer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

func TestInitiate_GetCbsStatusFailed(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, accountSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{}, errors.New("mock error"))

	res, err := uc.Initiate(nil, &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestInitiate_CbsNotReady(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, accountSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      true,
			IsStandIn:  false,
		}, nil)

	res, err := uc.Initiate(nil, &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestInitiate_FailedGetAccount(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, accountSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{}, errors.New("mock error"))

	res, err := uc.Initiate(nil, &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}
