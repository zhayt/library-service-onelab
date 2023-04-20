package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/transaction-service/model"
	"go.uber.org/zap"
)

type TransactionStorage struct {
	db *sqlx.DB
	l  *zap.Logger
}

func NewTransactionStorage(db *sqlx.DB, l *zap.Logger) *TransactionStorage {
	return &TransactionStorage{db: db, l: l}
}

func (r *TransactionStorage) CreateTransaction(ctx context.Context, transaction model.Transaction) (int, error) {
	qr := `INSERT INTO transaction (user_name, amount) VALUES ($1, $2)`

	var transactionID int64

	if err := r.db.GetContext(ctx, &transactionID, qr, transaction.UserName, transaction.Amount); err != nil {
		return 0, fmt.Errorf("couldn't create transaction: %w", err)
	}

	return int(transactionID), nil
}

func (r *TransactionStorage) CreateTransactionItem(ctx context.Context, item model.TransactionItem) (int, error) {
	qr := `INSERT INTO transaction_time (transaction_id, book_title, book_author, price) VALUES ($1, $2, $3, $4) RETURNING id`

	var itemID int64

	if err := r.db.GetContext(ctx, &itemID, qr, item.TransactionID, item.Book.Title, item.Book.Author, item.Book.Price); err != nil {
		return 0, fmt.Errorf("couldn't create transaction item: %w", err)
	}

	return int(itemID), nil
}

func (r *TransactionStorage) DeleteTransaction(ctx context.Context, transactionID int) error {
	qr := `DELETE FROM transaction WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, qr, transactionID); err != nil {
		return fmt.Errorf("couldn't delete transaction: %w", err)
	}

	return nil
}
