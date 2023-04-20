package handler

import (
	"github.com/zhayt/transaction-service/service"
	"go.uber.org/zap"
)

type Handler struct {
	transaction service.ITransactionService
	l           *zap.Logger
}

func NewHandler(service *service.Service, l *zap.Logger) *Handler {
	return &Handler{transaction: service.Transaction, l: l}
}
