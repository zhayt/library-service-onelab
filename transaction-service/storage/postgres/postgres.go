package postgres

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/transaction-service/config"
	"time"
)

const _timeOut = 5 * time.Second

func Dial(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg.DataBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't get pool connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), _timeOut)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("couldn't connect db: %w", err)
	}

	return db, nil
}
