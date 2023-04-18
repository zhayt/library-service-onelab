package service

import (
	"context"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
)

type IBIHistoryStorage interface {
	GetCurrentBorrowedBooks(ctx context.Context) ([]model.BorrowedBooks, error)
	GetBIHistoryLastMonth(ctx context.Context) ([]model.BorrowedBooks, error)
	CreateBIHistory(ctx context.Context, bIHistory model.BIHistory) error
	UpdateBIHistory(ctx context.Context, bIHistoryID int) (int, error)
	DeleteBIHistory(ctx context.Context, bIHistoryID int) error
}

type BIHistory struct {
	history IBIHistoryService
	log     *zap.Logger
}

func NewBIHistory(log *zap.Logger, history IBIHistoryService) *BIHistory {
	return &BIHistory{history: history, log: log}
}

func (s *BIHistory) CreateBIHistory(ctx context.Context, history model.BIHistory) error {
	// нужно ли проверять существует ли книга с таким ID и пользовотель,
	// если да то в каком слое?
	// в слое Handler, или в этом же сервисе могу вызвать метод сторадже который дастаем мне юзера

	return s.history.CreateBIHistory(ctx, history)
}

func (s *BIHistory) GetCurrentBorrowedBooks(ctx context.Context) ([]model.BorrowedBooks, error) {
	return s.history.GetCurrentBorrowedBooks(ctx)
}

func (s *BIHistory) GetBIHistoryLastMonth(ctx context.Context) ([]model.BorrowedBooks, error) {
	return s.history.GetBIHistoryLastMonth(ctx)
}

func (s *BIHistory) UpdateBIHistory(ctx context.Context, bIHistoryID int) (int, error) {
	return s.history.UpdateBIHistory(ctx, bIHistoryID)
}

func (s *BIHistory) DeleteBIHistory(ctx context.Context, bIHistoryID int) error {
	return s.history.DeleteBIHistory(ctx, bIHistoryID)
}
