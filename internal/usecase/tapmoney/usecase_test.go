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
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

func TestInquiry_Success(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       1000000,
			AccountNumber: "123",
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
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestInquiry_GetCbsFailed(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	log.Configure("test")

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

	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}

func TestInquiry_CbsNotReady(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	log.Configure("development")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2030-01-01",
			IsEOD:      true,
			IsStandIn:  false,
		}, nil)

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "123",
		Amount:        10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}

func TestInquiry_GetAccountFailed(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{}, errors.New("account not found"))

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "123",
		Amount:        10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestInquiry_AccountInsufficientBalance(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       5000,
			AccountNumber: "123",
		}, nil)

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "321",
		Amount:        10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.BadRequest().SetMsg("Insufficient balance"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestInquiry_FailedToInquiryPayment(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       5000000,
			AccountNumber: "123",
		}, nil)

	paymentSvc.EXPECT().Inquiry(mock.Anything, mock.Anything, mock.Anything).
		Return(payment.Payment{}, errors.New("inquiry failed"))

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "321",
		Amount:        10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.BadRequest().SetMsg("Inquiry failed"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestInquiry_FailedCreateTransaction(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       1000000,
			AccountNumber: "123",
		}, nil)
	paymentSvc.EXPECT().Inquiry(mock.Anything, mock.Anything, mock.Anything).
		Return(payment.Payment{
			ID: "pay-123",
		}, nil)
	txRepo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(errors.New("failed to create transaction"))

	resp, err := uc.Inquiry(context.Background(), &InquiryRequest{
		CardNumber:    "6013501000500719",
		SourceAccount: "123",
		Amount:        10000,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestPayment_Success(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	txRepo.EXPECT().GetByUUID(mock.Anything, mock.Anything).
		Return(transaction.Transaction{
			UUID:               "trx-123",
			SourceAccount:      "001201001479315",
			DestinationAccount: "6013501000500719",
			Amount:             10000,
			Status:             "pending",
			Note:               "test",
			Fee:                1500,
		}, nil)
	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
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
	assert.Equal(t, "inq-success", resp.Status)
	assert.Equal(t, "Payment successful", resp.Message)
	assert.Equal(t, int64(10000), resp.Amount)
	assert.Equal(t, int64(1500), resp.Fee)
	assert.Equal(t, "test", resp.Notes)

	t.Log(resp)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestPayment_GetCbsFailed(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	log.Configure("test")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{}, errors.New("get cbs status failed"))

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}

func TestPayment_CbsNotReady(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	log.Configure("development")

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2030-01-01",
			IsEOD:      true,
			IsStandIn:  false,
		}, nil)

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
}

func TestPayment_TransactionNotFound(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	txRepo.EXPECT().GetByUUID(mock.Anything, mock.Anything).
		Return(transaction.Transaction{}, errors.New("transaction not found"))

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.NotFound().SetMsg("Transaction was not found"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestPayment_TransactionAlreadyProcessed(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)
	txRepo.EXPECT().GetByUUID(mock.Anything, mock.Anything).
		Return(transaction.Transaction{
			UUID:               "trx-123",
			SourceAccount:      "001201001479315",
			DestinationAccount: "6013501000500719",
			Amount:             10000,
			Status:             "inq-success",
			Note:               "test",
			Fee:                1500,
		}, nil)

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.BadRequest().SetMsg("Transaction is already processed"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestPayment_FailedToGetSourceAccount(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, mock.Anything).
		Return(transaction.Transaction{
			UUID:               "trx-123",
			SourceAccount:      "001201001479315",
			DestinationAccount: "6013501000500719",
			Amount:             10000,
			Status:             "pending",
			Note:               "test",
			Fee:                1500,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{}, errors.New("account not found"))

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.NotFound().SetMsg("Source account was not found"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestPayment_InsufficientBalance(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, mock.Anything).
		Return(transaction.Transaction{
			UUID:               "trx-123",
			SourceAccount:      "001201001479315",
			DestinationAccount: "6013501000500719",
			Amount:             10000,
			Status:             "pending",
			Note:               "test",
			Fee:                1500,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       5000,
			AccountNumber: "001201001479315",
		}, nil)

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.BadRequest().SetMsg("Insufficient balance"), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestPayment_FailedToProcessPayment(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, mock.Anything).
		Return(transaction.Transaction{
			UUID:               "trx-123",
			SourceAccount:      "001201001479315",
			DestinationAccount: "6013501000500719",
			Amount:             10000,
			Status:             "pending",
			Note:               "test",
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       1000000,
			AccountNumber: "001201001479315",
		}, nil)

	paymentSvc.EXPECT().Payment(mock.Anything, mock.Anything).
		Return(payment.Payment{}, errors.New("payment failed"))

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}

func TestPayment_FailedToUpdateTransaction(t *testing.T) {
	var (
		cbsService = cbs.NewMockService(t)
		txRepo     = transaction.NewMockRepository(t)
		paymentSvc = payment.NewMockService(t)
		accountSvc = account.NewMockService(t)
		uc         = NewUsecase(cbsService, txRepo, paymentSvc, accountSvc)
	)

	cbsService.EXPECT().GetStatus(mock.Anything).
		Return(cbs.Status{
			SystemDate: "2025-08-21",
			IsEOD:      false,
			IsStandIn:  false,
		}, nil)

	txRepo.EXPECT().GetByUUID(mock.Anything, mock.Anything).
		Return(transaction.Transaction{
			UUID:               "trx-123",
			SourceAccount:      "001201001479315",
			DestinationAccount: "6013501000500719",
			Amount:             10000,
			Status:             "pending",
			Note:               "test",
			Fee:                1500,
		}, nil)

	accountSvc.EXPECT().Get(mock.Anything, mock.Anything).
		Return(account.Account{
			Balance:       1000000,
			AccountNumber: "001201001479315",
		}, nil)

	paymentSvc.EXPECT().Payment(mock.Anything, mock.Anything).
		Return(payment.Payment{
			Status: "success",
		}, nil)

	txRepo.EXPECT().Update(mock.Anything, mock.Anything).
		Return(errors.New("failed to update transaction"))

	resp, err := uc.Payment(context.Background(), &PaymentRequest{
		TransactionID: "trx-123",
		Amount:        10000,
		Notes:         "test",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, pkgerror.InternalServerError(), err)

	cbsService.AssertExpectations(t)
	txRepo.AssertExpectations(t)
	paymentSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}
