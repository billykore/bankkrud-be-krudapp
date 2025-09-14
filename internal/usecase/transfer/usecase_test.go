package transfer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transfer"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

func TestInitiate_GetCbsStatusFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountSvc  = account.NewMockService(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountSvc, transferSvc)
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
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountSvc  = account.NewMockService(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountSvc, transferSvc)
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

func TestInitiate_GetSourceAccountFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountSvc  = account.NewMockService(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountSvc, transferSvc)
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

func TestInitiate_SourceAccountCannotTransfer(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountSvc  = account.NewMockService(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountSvc, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       5000,
		}, nil)

	res, err := uc.Initiate(nil, &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.BadRequest().SetMsg("Insufficient balance"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestInitiate_GetDestinationAccountFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountSvc  = account.NewMockService(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountSvc, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, "123").
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       50000,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, "456").
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

func TestInitiate_CreateTransactionFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountSvc  = account.NewMockService(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountSvc, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, "123").
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       50000,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, "456").
		Return(account.Account{
			AccountNumber: "456",
			FullName:      "Jane Doe",
			Type:          "savings",
			Balance:       30000,
		}, nil)

	txRepo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(errors.New("mock error"))

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

func TestInitiate_Success(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountSvc  = account.NewMockService(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountSvc, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, "123").
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       50000,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, "456").
		Return(account.Account{
			AccountNumber: "456",
			FullName:      "Jane Doe",
			Type:          "savings",
			Balance:       30000,
		}, nil)

	txRepo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(nil)

	res, err := uc.Initiate(nil, &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.NotNil(t, res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.TransactionID)
	assert.Equal(t, "inq-success", res.Status)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}
