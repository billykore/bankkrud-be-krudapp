package repo

import (
	"context"
	"errors"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/storage/model"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
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
	var m model.Transaction
	res := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&m)
	if res.Error != nil {
		return transaction.Transaction{}, res.Error
	}
	return transaction.Transaction{
		UUID:                 m.UUID,
		SourceAccount:        m.SourceAccount,
		DestinationAccount:   m.DestinationAccount,
		TransactionType:      m.TransactionType,
		TransactionReference: m.TransactionReference,
		Status:               m.Status,
		Notes:                m.Note,
		Amount:               m.Amount,
		Fee:                  m.Fee,
	}, nil
}

func (r *TransactionRepo) Create(ctx context.Context, tx transaction.Transaction) error {
	res := r.db.WithContext(ctx).Create(&model.Transaction{
		UUID:                 tx.UUID,
		SourceAccount:        tx.SourceAccount,
		DestinationAccount:   tx.DestinationAccount,
		TransactionType:      tx.TransactionType,
		TransactionReference: tx.TransactionReference,
		Status:               tx.Status,
		Note:                 tx.Notes,
		Amount:               tx.Amount,
		Fee:                  tx.Fee,
	})
	return res.Error
}

func (r *TransactionRepo) Update(ctx context.Context, tx transaction.Transaction) error {
	return errors.New("not implemented")
}
