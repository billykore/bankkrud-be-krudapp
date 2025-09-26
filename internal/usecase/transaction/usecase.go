package transaction

import (
	"context"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

type Usecase struct {
	txRepo transaction.Repository
}

func NewUsecase(txRepo transaction.Repository) *Usecase {
	return &Usecase{
		txRepo: txRepo,
	}
}

func (uc *Usecase) GetTransactions(ctx context.Context, req *GetTransactionsRequest) ([]*TransactionDataResponse, error) {
	l := log.WithContext(ctx, "Detail")

	userFromCtx, err := user.FromContext(ctx)
	if err != nil {
		l.Error().Err(err).Msg("Error getting user from context")
		return nil, pkgerror.Unauthorized().SetMsg("User not authorized")
	}

	requestMap := req.Map()
	if userFromCtx.Username != "" {
		requestMap["user_username"] = userFromCtx.Username
	}

	txs, err := uc.txRepo.GetByParams(ctx, requestMap)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get transactions")
		return nil, pkgerror.NotFound().SetMsg("Failed to get transactions")
	}

	details := make([]*TransactionDataResponse, 0, 10)
	for _, tx := range txs {
		details = append(details, &TransactionDataResponse{
			UUID:               tx.UUID,
			TransactionType:    tx.TransactionType,
			Status:             tx.Status,
			SourceAccount:      tx.SourceAccount,
			DestinationAccount: tx.DestinationAccount,
			Amount:             tx.Amount,
			ProcessedAt:        tx.ProcessedAt,
		})
	}

	return details, nil
}

func (uc *Usecase) GetTransactionByUUID(ctx context.Context, req *GetTransactionRequest) (*TransactionDataResponse, error) {
	l := log.WithContext(ctx, "Detail")

	tx, err := uc.txRepo.GetByUUID(ctx, req.UUID)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get transaction")
		return nil, pkgerror.NotFound().SetMsg("Failed to get transaction")
	}

	return &TransactionDataResponse{
		UUID:               tx.UUID,
		TransactionType:    tx.TransactionType,
		Status:             tx.Status,
		SourceAccount:      tx.SourceAccount,
		DestinationAccount: tx.DestinationAccount,
		Amount:             tx.Amount,
		ProcessedAt:        tx.ProcessedAt,
	}, nil
}
