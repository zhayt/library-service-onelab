package service

import (
	"context"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/storage"
	"go.uber.org/zap"
)

type IBookStorage interface {
	GetBookById(ctx context.Context, bookID int) (model.Book, error)
	GetAllBooks(ctx context.Context) ([]model.Book, error)
	CreateBook(ctx context.Context, book model.Book) (int, error)
	UpdateBook(ctx context.Context, book model.Book) (int, error)
	DeleteBook(ctx context.Context, bookID int) error
}

type BookService struct {
	book storage.IBookStorage
	log  *zap.Logger
}

func (s *BookService) CreateBook(ctx context.Context, book model.Book) (model.Book, error) {
	// проверка валидности данных
	bookId, err := s.book.CreateBook(ctx, book)
	if err != nil {
		return book, err
	}

	book.ID = bookId

	return book, nil
}

func (s *BookService) GetBookById(ctx context.Context, bookId int) (model.Book, error) {
	return s.book.GetBookById(ctx, bookId)
}

func (s *BookService) GetBooks(ctx context.Context) ([]model.Book, error) {
	return s.book.GetAllBooks(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, book model.Book) (int, error) {
	// проверка валидности данных

	return s.book.UpdateBook(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, bookId int) error {
	return s.book.DeleteBook(ctx, bookId)
}

func NewBookService(log *zap.Logger, book storage.IBookStorage) *BookService {
	return &BookService{book: book, log: log}
}
