package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
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

func (r *BookStorage) GetBookByID(ctx context.Context, bookID int) (model.Book, error) {
	qr := `SELECT * FROM book WHERE id = $1`

	var book model.Book
	if err := r.db.GetContext(ctx, &book, qr, bookID); err != nil {
		return book, fmt.Errorf("couldn't take book id#%v: %w", bookID, err)
	}

	return book, nil
}

func (r *BookStorage) GetAllBooks(ctx context.Context) ([]model.Book, error) {
	qr := `SELECT * FROM book`

	var books []model.Book

	if err := r.db.SelectContext(ctx, &books, qr); err != nil {
		return books, fmt.Errorf("couldn't take all books: %w", err)
	}

	return books, nil
}

func (r *BookStorage) CreateBook(ctx context.Context, book model.Book) (int, error) {
	qr := `INSERT INTO book (name, author) VALUES($1, $2) RETURNING id`

	var bookID int64
	if err := r.db.GetContext(ctx, &bookID, qr, book.Name, book.Author); err != nil {
		return 0, fmt.Errorf("couldn't create book: %w", err)
	}

	return int(bookID), nil
}

func (r *BookStorage) UpdateBook(ctx context.Context, book model.Book) (int, error) {
	qr := `UPDATE book SET name = $2, author = $3 WHERE id = $1 RETURNING id`

	var bookId int64

	if err := r.db.GetContext(ctx, &bookId, qr, book.ID, book.Name, book.Author); err != nil {
		return 0, fmt.Errorf("couldn't update book id#%v: %w", book.ID, err)
	}

	return int(bookId), nil
}

func (r *BookStorage) DeleteBook(ctx context.Context, bookID int) error {
	qr := `DELETE FROM book WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, qr, bookID); err != nil {
		return fmt.Errorf("cannot delete book id#%v: %w", bookID, err)
	}

	return nil
}
