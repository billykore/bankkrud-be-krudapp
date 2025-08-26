package repo

import (
	"context"
	"errors"

	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/storage/model"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/transaction"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

func (r *TransactionRepo) Get(ctx context.Context, uuid string) (transaction.Transaction, error) {
	return transaction.Transaction{}, errors.New("not implemented")
}

func (r *TransactionRepo) Create(ctx context.Context, tx transaction.Transaction) error {
	return errors.New("not implemented")
}

func (r *TransactionRepo) Update(ctx context.Context, tx transaction.Transaction) error {
	res := r.db.WithContext(ctx).Where(`"UUID" = ?`, tx.UUID).
		Updates(&model.Transaction{
			Status: tx.Status,
		})
	return res.Error
}
