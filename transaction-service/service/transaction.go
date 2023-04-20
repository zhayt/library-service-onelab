package service

import (
	"context"
	"github.com/zhayt/transaction-service/model"
	"github.com/zhayt/transaction-service/storage"
	"go.uber.org/zap"
)

type TransactionService struct {
	storage storage.ITransactionStorage
	l       *zap.Logger
}

func NewTransactionService(storage storage.ITransactionStorage, l *zap.Logger) *TransactionService {
	return &TransactionService{storage: storage, l: l}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction model.Transaction) (int, error) {
	return s.storage.CreateTransaction(ctx, transaction)
}

func (s *TransactionService) CreateTransactionItem(ctx context.Context, item model.TransactionItem) (int, error) {
	return s.storage.CreateTransactionItem(ctx, item)
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, transactionID int) error {
	return s.storage.DeleteTransaction(ctx, transactionID)
}
