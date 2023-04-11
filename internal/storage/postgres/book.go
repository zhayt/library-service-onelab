package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
)

type BookStorage struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewBookStorage(db *sqlx.DB, logger *zap.Logger) *BookStorage {
	return &BookStorage{db: db, log: logger}
}

func (r *BookStorage) GetBookById(ctx context.Context, bookID int) (model.Book, error) {
	qr := `SELECT * FROM book WHERE id = $1`

	var book model.Book
	if err := r.db.GetContext(ctx, &book, qr); err != nil {
		return book, fmt.Errorf("cannot take book: %w", common.ErrBookNotExists)
	}

	return book, nil
}

func (r *BookStorage) GetAllBooks(ctx context.Context) ([]model.Book, error) {
	qr := `SELECT * FROM book`

	var books []model.Book

	if err := r.db.SelectContext(ctx, &books, qr); err != nil {
		return books, fmt.Errorf("cannot take all books: %w", common.ErrNoRows)
	}

	return books, nil
}

func (r *BookStorage) CreateBook(ctx context.Context, book model.Book) (int, error) {
	qr := `INSERT INTO book (name, author) VALUES($1, $2) RETURNING id`

	var bookID int64
	if err := r.db.GetContext(ctx, &bookID, qr, book.Name, book.Author); err != nil {
		return 0, fmt.Errorf("cannot create book: %w", err)
	}

	return int(bookID), nil
}

func (r *BookStorage) UpdateBook(ctx context.Context, book model.Book) (int, error) {
	qr := `UPDATE book SET name = $2, author = $3 WHERE id = $1 RETURNING id`

	var bookId int64

	if err := r.db.GetContext(ctx, &bookId, qr, book.ID, book.Name, book.Author); err != nil {
		return 0, fmt.Errorf("cannot update book: %w", err)
	}

	return int(bookId), nil
}

func (r *BookStorage) DeleteBook(ctx context.Context, bookID int) error {
	qr := `DELETE FROM book WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, qr, bookID); err != nil {
		return fmt.Errorf("cannot delete book: %w", err)
	}

	return nil
}
