package transfer

import (
	"context"
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
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{}, errors.New("mock error"))

	res, err := uc.Initiate(context.Background(), &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestInitiate_CbsNotReady(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      true,
			IsStandIn:  false,
		}, nil)

	res, err := uc.Initiate(context.Background(), &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestInitiate_GetSourceAccountFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{}, errors.New("mock error"))

	res, err := uc.Initiate(context.Background(), &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestInitiate_SourceAccountCannotTransfer(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       5000,
		}, nil)

	res, err := uc.Initiate(context.Background(), &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.BadRequest().SetMsg("Insufficient balance"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestInitiate_GetDestinationAccountFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, "123").
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       50000,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, "456").
		Return(account.Account{}, errors.New("mock error"))

	res, err := uc.Initiate(context.Background(), &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestInitiate_CreateTransactionFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, "123").
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       50000,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, "456").
		Return(account.Account{
			AccountNumber: "456",
			FullName:      "Jane Doe",
			Type:          "savings",
			Balance:       30000,
		}, nil)

	txRepo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(errors.New("mock error"))

	res, err := uc.Initiate(context.Background(), &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestInitiate_Success(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, "123").
		Return(account.Account{
			AccountNumber: "123",
			FullName:      "John Doe",
			Type:          "savings",
			Balance:       50000,
		}, nil)

	accountRepo.EXPECT().Get(mock.Anything, "456").
		Return(account.Account{
			AccountNumber: "456",
			FullName:      "Jane Doe",
			Type:          "savings",
			Balance:       30000,
		}, nil)

	txRepo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(nil)

	res, err := uc.Initiate(context.Background(), &InitiateRequest{
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.NotNil(t, res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.UUID)
	assert.Equal(t, "inq-success", res.Status)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestProcess_GetCbsStatusFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{}, errors.New("mock error"))

	res, err := uc.Process(context.Background(), &ProcessRequest{
		UUID:               "tx-123",
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestProcess_CbsNotReady(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      true,
			IsStandIn:  false,
		}, nil)

	res, err := uc.Process(context.Background(), &ProcessRequest{
		UUID:               "tx-123",
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestProcess_GetTransactionFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, "tx-123").
		Return(transaction.Transaction{}, errors.New("mock error"))

	res, err := uc.Process(context.Background(), &ProcessRequest{
		UUID:               "tx-123",
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestProcess_TransactionStatusNotInquirySuccess(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, "tx-123").
		Return(transaction.Transaction{
			UUID:   "tx-123",
			Status: "completed",
		}, nil)

	res, err := uc.Process(context.Background(), &ProcessRequest{
		UUID:               "tx-123",
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t,
		pkgerror.BadRequest().SetMsg("Transaction is not in a valid state to be processed"),
		err,
	)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestProcess_TransferFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, "tx-123").
		Return(transaction.Transaction{
			UUID:               "tx-123",
			Status:             transaction.StatusCompleted,
			SourceAccount:      "123",
			DestinationAccount: "456",
		}, nil)

	transferSvc.EXPECT().Transfer(
		mock.Anything,
		"123",
		"456",
		int64(10000),
		"TRF 123 456 BNKKRD tx-123",
	).Return(transfer.Transfer{}, errors.New("mock error"))

	res, err := uc.Process(context.Background(), &ProcessRequest{
		UUID:               "tx-123",
		SourceAccount:      "123",
		DestinationAccount: "456",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
	transferSvc.AssertExpectations(t)
}

func TestProcess_UpdateTransactionFailed(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, "tx-123").
		Return(transaction.Transaction{
			UUID:               "tx-123",
			Status:             transaction.StatusCompleted,
			SourceAccount:      "121",
			DestinationAccount: "454",
		}, nil)

	transferSvc.EXPECT().Transfer(
		mock.Anything,
		"121",
		"454",
		int64(10000),
		"TRF 121 454 BNKKRD tx-123",
	).Return(transfer.Transfer{
		TransactionReference: "ref-123",
		Status:               "success",
	}, nil)

	txRepo.EXPECT().Update(mock.Anything, mock.Anything).
		Return(errors.New("mock error"))

	res, err := uc.Process(context.Background(), &ProcessRequest{
		UUID:               "tx-123",
		SourceAccount:      "121",
		DestinationAccount: "45a",
		Amount:             10000,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
	transferSvc.AssertExpectations(t)
}

func TestProcess_Success(t *testing.T) {
	var (
		cbsService  = cbs.NewMockService(t)
		txRepo      = transaction.NewMockRepository(t)
		accountRepo = account.NewMockRepository(t)
		transferSvc = transfer.NewMockService(t)
		uc          = NewUsecase(cbsService, txRepo, accountRepo, transferSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, "tx-123").
		Return(transaction.Transaction{
			UUID:               "tx-123",
			Status:             transaction.StatusCompleted,
			SourceAccount:      "121",
			DestinationAccount: "454",
		}, nil)

	transferSvc.EXPECT().Transfer(
		mock.Anything,
		"121",
		"454",
		int64(10000),
		"TRF 121 454 BNKKRD tx-123",
	).Return(transfer.Transfer{
		TransactionReference: "ref-123",
		Status:               "success",
	}, nil)

	txRepo.EXPECT().Update(mock.Anything, mock.Anything).
		Return(nil)

	res, err := uc.Process(context.Background(), &ProcessRequest{
		UUID:               "tx-123",
		SourceAccount:      "121",
		DestinationAccount: "454",
		Amount:             10000,
	})

	assert.NotNil(t, res)
	assert.NoError(t, err)
	assert.Equal(t, "tx-123", res.UUID)
	assert.Equal(t, "success", res.Status)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
	transferSvc.AssertExpectations(t)
}
