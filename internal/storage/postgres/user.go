package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
)

type UserStorage struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewUserStorage(db *sqlx.DB, logger *zap.Logger) *UserStorage {
	return &UserStorage{db: db, log: logger}
}

func (r *UserStorage) GetUserById(ctx context.Context, userID int) (model.User, error) {
	qr := `SELECT * FROM "user" WHERE id = $1 LIMIT 1`

	var user model.User

	if err := r.db.GetContext(ctx, &user, qr, userID); err != nil {
		r.log.Error("Storage: GetUserById error", zap.Error(err))
		return user, err
	}

	return user, nil
}

func (r *UserStorage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	qr := `SELECT * FROM "user" WHERE email = $1`

	var user model.User

	if err := r.db.GetContext(ctx, &user, qr, email); err != nil {
		r.log.Error("Storage: GetUserByEmail error", zap.Error(err))
		return user, err
	}

	return user, nil
}

func (r *UserStorage) CreateUser(ctx context.Context, user model.User) (int, error) {
	qr := `INSERT INTO "user" (fio, email, password) VALUES ($1, $2, $3) RETURNING id`

	var userID int64

	if err := r.db.GetContext(ctx, &userID, qr, user.FIO, user.Email, user.Password); err != nil {
		log.Error("Storage create user error", zap.Error(err))
		return 0, fmt.Errorf("cannot create user: %w", err)
	}

	return int(userID), nil
}

func (r *UserStorage) UpdateUserFIO(ctx context.Context, user model.UserUpdateFIO) (int, error) {
	qr := `UPDATE "user" SET fio = $2 WHERE id = $1`

	var userID int64
	if err := r.db.GetContext(ctx, &userID, qr, user.ID, user.FIO); err != nil {
		return 0, fmt.Errorf("cannot update user: %w", err)
	}

	return int(userID), nil
}

func (r *UserStorage) UpdateUserPassword(ctx context.Context, user model.UserUpdatePassword) (int, error) {
	qr := `UPDATE "user" SET password = $2 WHERE id = $1`

	var userID int64
	if err := r.db.GetContext(ctx, &userID, qr, user.ID, user.NewPassword); err != nil {
		return 0, fmt.Errorf("cannot update user: %w", err)
	}

	return int(userID), nil
}

func (r *UserStorage) DeleteUser(ctx context.Context, userID int) error {
	qr := `DELETE FROM "user" WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, qr, userID); err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
