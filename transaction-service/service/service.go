package service

import (
	"context"
	"github.com/zhayt/transaction-service/model"
	"github.com/zhayt/transaction-service/storage"
	"go.uber.org/zap"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction) (int, error)
	CreateTransactionItem(ctx context.Context, item model.TransactionItem) (int, error)
	DeleteTransaction(ctx context.Context, transactionID int) error
}

type Service struct {
	Transaction ITransactionService
}

func NewService(storage *storage.Storage, logger *zap.Logger) *Service {
	return &Service{Transaction: NewTransactionService(storage, logger)}
}
