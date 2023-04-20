package storage

import (
	"context"
	"github.com/zhayt/transaction-service/config"
	"github.com/zhayt/transaction-service/model"
	"github.com/zhayt/transaction-service/storage/postgres"
	"go.uber.org/zap"
)

type ITransactionStorage interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction) (int, error)
	CreateTransactionItem(ctx context.Context, item model.TransactionItem) (int, error)
	DeleteTransaction(ctx context.Context, transactionID int) error
}

type Storage struct {
	ITransactionStorage
}

func NewStorage(cfg *config.Config, logger *zap.Logger) (*Storage, error) {
	db, err := postgres.Dial(cfg)
	if err != nil {
		return nil, err
	}

	return &Storage{ITransactionStorage: postgres.NewTransactionStorage(db, logger)}, nil
}
