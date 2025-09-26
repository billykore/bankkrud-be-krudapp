package repo

import (
	"context"

	"github.com/google/uuid"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/storage/model"
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

func (r *TransactionRepo) GetByUUID(ctx context.Context, tfuuid string) (transaction.Transaction, error) {
	var m model.Transaction
	id, err := uuid.Parse(tfuuid)
	if err != nil {
		return transaction.Transaction{}, err
	}
	res := r.db.WithContext(ctx).
		Where("uuid = ?", id).
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
		Note:                 m.Note,
		Amount:               m.Amount,
		Fee:                  m.Fee,
		ProcessedAt:          m.CreatedAt,
	}, nil
}

func (r *TransactionRepo) GetByParams(ctx context.Context, params map[string]any) ([]transaction.Transaction, error) {
	var models []model.Transaction
	res := r.db.WithContext(ctx).
		Where(params).
		Find(&models)
	if res.Error != nil {
		return nil, res.Error
	}
	var transactions []transaction.Transaction
	for _, m := range models {
		transactions = append(transactions, transaction.Transaction{
			UUID:                 m.UUID,
			SourceAccount:        m.SourceAccount,
			DestinationAccount:   m.DestinationAccount,
			TransactionType:      m.TransactionType,
			TransactionReference: m.TransactionReference,
			Status:               m.Status,
			Note:                 m.Note,
			Amount:               m.Amount,
			Fee:                  m.Fee,
			ProcessedAt:          m.CreatedAt,
		})
	}
	return transactions, nil
}

func (r *TransactionRepo) Create(ctx context.Context, tx transaction.Transaction) error {
	res := r.db.WithContext(ctx).Create(&model.Transaction{
		UUID:                 tx.UUID,
		SourceAccount:        tx.SourceAccount,
		DestinationAccount:   tx.DestinationAccount,
		TransactionType:      tx.TransactionType,
		TransactionReference: tx.TransactionReference,
		Status:               tx.Status,
		Note:                 tx.Note,
		Amount:               tx.Amount,
		Fee:                  tx.Fee,
		UserUsername:         tx.Username,
	})
	return res.Error
}

func (r *TransactionRepo) Update(ctx context.Context, tx transaction.Transaction) error {
	res := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("uuid = ?", tx.UUID).
		Updates(&model.Transaction{
			SourceAccount:        tx.SourceAccount,
			DestinationAccount:   tx.DestinationAccount,
			TransactionType:      tx.TransactionType,
			TransactionReference: tx.TransactionReference,
			Status:               tx.Status,
			Note:                 tx.Note,
			Amount:               tx.Amount,
			Fee:                  tx.Fee,
		})
	return res.Error
}
