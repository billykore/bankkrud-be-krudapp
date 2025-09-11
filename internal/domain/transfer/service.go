package transfer

import "context"

type Service interface {
	// Transfer moves amount from one account to another
	// and returns an error if the operation fails.
	Transfer(ctx context.Context, srcAccountNumber, destAccountNumber string, amount int64, remark string) (Transfer, error)
}
