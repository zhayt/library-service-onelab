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

func (r *BIHistoryStorage) GetCurrentBorrowedBooks(ctx context.Context) ([]model.BIHistory, error) {
	qr := `SELECT * FROM book_issue_history WHERE return_date IS NULL`

	var bIHistories []model.BIHistory

	if err := r.db.SelectContext(ctx, &bIHistories, qr); err != nil {
		return nil, fmt.Errorf("couldn't teke book issue history: %w", err)
	}

	return bIHistories, nil
}

func (r *BIHistoryStorage) GetBIHistoryLastMonth(ctx context.Context) ([]model.BIHistory, error) {
	qr := `SELECT * FROM book_issue_history
		WHERE issue_date >= NOW() - INTERVAL '1 month' AND return_date IS NULL;`

	var bIHistories []model.BIHistory

	if err := r.db.SelectContext(ctx, &bIHistories, qr); err != nil {
		return bIHistories, fmt.Errorf("couldn't take book issue history for last month: %w", err)
	}

	return bIHistories, nil
}

func (r *BIHistoryStorage) CreateBIHistory(ctx context.Context, bIHistory model.BIHistory) (int, error) {
	qr := `INSERT INTO book_issue_history (book_id, user_id, issue_date) VALUES ($1, $2, $3) RETURNING id`

	var bihId int64

	err := r.db.GetContext(ctx, &bihId, qr, bIHistory.BookID, bIHistory.UserID, bIHistory.IssueDate)
	if err != nil {
		return 0, fmt.Errorf("couldn't create book issue history: %w", err)
	}
	return int(bihId), nil
}

func (r *BIHistoryStorage) UpdateBIHistory(ctx context.Context, bIHistory model.BIHistory) (int, error) {
	qr := `UPDATE book_issue_history SET return_date = $2 WHERE id = $1 RETURNING id`

	var bihId int64

	if err := r.db.GetContext(ctx, &bihId, qr, bIHistory.ID, bIHistory.ReturnDate); err != nil {
		return 0, fmt.Errorf("couldn't update book issue history returning date: %w", err)
	}

	return int(bihId), nil
}

func (r *BIHistoryStorage) DeleteBIHistory(ctx context.Context, bIHistoryID int) error {
	qr := `DELETE FROM book_issue_history WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, qr, bIHistoryID); err != nil {
		return fmt.Errorf("couldn't delete book ID#%v: %w", bIHistoryID, err)
	}

	return nil
}
