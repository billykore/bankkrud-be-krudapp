package api

import (
	"context"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transfer"
)

type CBSTransferAPI struct{}

func NewCBSTransferAPI() *CBSTransferAPI {
	return &CBSTransferAPI{}
}

func (ta *CBSTransferAPI) Transfer(ctx context.Context, srcAccountNumber, destAccountNumber string, amount int64, remark string) (transfer.Transfer, error) {
	return transfer.Transfer{
		SourceAccount:        srcAccountNumber,
		DestinationAccount:   destAccountNumber,
		Amount:               amount,
		Status:               "success",
		TransactionReference: "example-ref-123",
	}, nil
}
