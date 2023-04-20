package service

import (
	"context"
	"fmt"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/storage"
	"go.uber.org/zap"
)

type IBIHistoryService interface {
	CreateBIHistory(ctx context.Context, history model.BIHistory) error
	GetCurrentBorrowedBooks(ctx context.Context) ([]model.BorrowedBooks, error)
	GetBIHistoryLastMonth(ctx context.Context) ([]model.BorrowedBooks, error)
	UpdateBIHistory(ctx context.Context, bIHistoryID int) (int, error)
	DeleteBIHistory(ctx context.Context, bIHistoryID int) error
}

type ITransactionService interface {
	CreateTransaction(transaction model.Transaction) (int, error)
	CreateTransactionItem(item model.TransactionItem) error
	DeleteTransaction(transactionID int) error
}

type IGetBookUser interface {
	GetBookByID(ctx context.Context, bookID int) (model.Book, error)
	GetUserByID(ctx context.Context, userID int) (model.User, error)
}

type RentTransactionService struct {
	IBIHistoryService
	ITransactionService
	IGetBookUser
	l *zap.Logger
}

func NewRentTransactionService(logger *zap.Logger, storage *storage.Storage) *RentTransactionService {
	return &RentTransactionService{IBIHistoryService: NewBIHistory(logger, storage), ITransactionService: NewTransaction(logger), IGetBookUser: storage}
}

func (s *RentTransactionService) RentBook(ctx context.Context, history model.BIHistory) error {
	if err := s.CreateBIHistory(ctx, history); err != nil {
		return fmt.Errorf("couldn't create bihistory: %w", err)
	}

	user, err := s.GetUserByID(ctx, history.UserID)
	if err != nil {
		return fmt.Errorf("couldn't create bihistory: %w", err)
	}

	books := make([]model.Book, 0, len(history.Books))

	var amount float64

	for _, rentBook := range history.Books {
		book, err := s.GetBookByID(ctx, rentBook.ID)
		if err != nil {
			return fmt.Errorf("couldn't create bihistory: %w", err)
		}

		amount += book.Price * float64(rentBook.Quantity)

		books = append(books, book)
	}

	transaction := model.Transaction{
		UserName: user.FIO,
		Amount:   amount,
	}

	transactionID, err := s.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("couldn't create bihistory: %w", err)
	}

	for _, book := range books {
		item := model.TransactionItem{
			TransactionID: uint(transactionID),
			Book:          &book,
		}

		if err = s.CreateTransactionItem(item); err != nil {
			if err = s.DeleteTransaction(transactionID); err != nil {
				return fmt.Errorf("delete transaction error")
			}

			// удалить записи из BIHistory
			return fmt.Errorf("couldn't create bihistory: %w", err)
		}
	}

	return nil
}
