package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
)

type BIHistoryStorage struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewBIHistory(db *sqlx.DB, logger *zap.Logger) *BIHistoryStorage {
	return &BIHistoryStorage{db: db, log: logger}
}

func (r *BIHistoryStorage) CreateBIHistory(ctx context.Context, bIHistory model.BIHistory) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("couldn't create book issue history: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO book_issue_history (book_id, quantity, user_id) 
		   VALUES ($1, $2, $3)`)

	if err != nil {
		return fmt.Errorf("couldn't prepare query: %w", err)
	}

	for _, book := range bIHistory.Books {
		if _, err := stmt.ExecContext(ctx, book.ID, book.Quantity, bIHistory.UserID); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("couldn't execute query: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return nil
}
func (r *BIHistoryStorage) GetCurrentBorrowedBooks(ctx context.Context) ([]model.BorrowedBooks, error) {
	qr := `SELECT book_issue_history.id, u.fio, b.name, b.author, quantity, created_at FROM 
           book_issue_history
		   INNER JOIN "user" u on u.id = book_issue_history.user_id
		   INNER JOIN book b on b.id = book_issue_history.book_id
           WHERE return_date IS NULL`

	var borrowedBooks []model.BorrowedBooks

	if err := r.db.SelectContext(ctx, &borrowedBooks, qr); err != nil {
		return nil, fmt.Errorf("couldn't teke book issue history: %w", err)
	}

	return borrowedBooks, nil
}

func (r *BIHistoryStorage) GetBIHistoryLastMonth(ctx context.Context) ([]model.BorrowedBooks, error) {
	qr := `SELECT book_issue_history.id, u.fio, b.name, b.author, created_at FROM 
           book_issue_history
		   INNER JOIN "user" u on u.id = book_issue_history.user_id
		   INNER JOIN book b on b.id = book_issue_history.book_id
		   WHERE created_at >= NOW() - INTERVAL '1 month' AND return_date IS NULL;`

	var bIHistories []model.BorrowedBooks

	if err := r.db.SelectContext(ctx, &bIHistories, qr); err != nil {
		return bIHistories, fmt.Errorf("couldn't take book issue history for last month: %w", err)
	}

	return bIHistories, nil
}

func (r *BIHistoryStorage) UpdateBIHistory(ctx context.Context, bIHistoryID int) (int, error) {
	qr := `UPDATE book_issue_history 
		   SET return_date = CURRENT_TIMESTAMP 
		   WHERE id = $1 
		   RETURNING id`

	var bihId int64

	if err := r.db.GetContext(ctx, &bihId, qr, bIHistoryID); err != nil {
		return 0, fmt.Errorf("couldn't update book issue history returning date: %w", err)
	}

	return int(bihId), nil
}

func (r *BIHistoryStorage) DeleteBIHistory(ctx context.Context, bIHistoryID int) error {
	qr := `DELETE FROM book_issue_history
       	   WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, qr, bIHistoryID); err != nil {
		return fmt.Errorf("couldn't delete book ID#%v: %w", bIHistoryID, err)
	}

	return nil
}
