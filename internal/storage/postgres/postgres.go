package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const _timeOut = 5 * time.Second

func Dial(driver string, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("connot open db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), _timeOut)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	return db, nil
}
