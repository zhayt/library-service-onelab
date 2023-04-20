package service

import (
	"context"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
)

type IBookStorage interface {
	GetBookByID(ctx context.Context, bookID int) (model.Book, error)
	GetAllBooks(ctx context.Context) ([]model.Book, error)
	CreateBook(ctx context.Context, book model.Book) (int, error)
	UpdateBook(ctx context.Context, book model.Book) (int, error)
	DeleteBook(ctx context.Context, bookID int) error
}

type BookService struct {
	book IBookStorage
	log  *zap.Logger
}

func NewBookService(log *zap.Logger, book IBookStorage) *BookService {
	return &BookService{book: book, log: log}
}

func (s *BookService) CreateBook(ctx context.Context, book model.Book) (int, error) {
	// проверка валидности данных
	return s.book.CreateBook(ctx, book)
}

func (s *BookService) GetBookByID(ctx context.Context, bookId int) (model.Book, error) {
	return s.book.GetBookByID(ctx, bookId)
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]model.Book, error) {
	return s.book.GetAllBooks(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, book model.Book) (int, error) {
	// проверка валидности данных

	return s.book.UpdateBook(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, bookId int) error {
	return s.book.DeleteBook(ctx, bookId)
}
