package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/storage/postgres"
	"go.uber.org/zap"
)

type IUserStorage interface {
	GetUserById(ctx context.Context, userID int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)

	CreateUser(ctx context.Context, user model.User) (int, error)

	UpdateUserFIO(ctx context.Context, user model.UserUpdateFIO) (int, error)
	UpdateUserPassword(ctx context.Context, user model.UserUpdatePassword) (int, error)

	DeleteUser(ctx context.Context, userID int) error
}

type IBookStorage interface {
	GetBookById(ctx context.Context, bookID int) (model.Book, error)
	GetAllBooks(ctx context.Context) ([]model.Book, error)
	CreateBook(ctx context.Context, book model.Book) (int, error)
	UpdateBook(ctx context.Context, book model.Book) (int, error)
	DeleteBook(ctx context.Context, bookID int) error
}

type IBIHistoryStorage interface {
	GetBIHistory(ctx context.Context) ([]model.BIHistory, error)
	GetBIHistoryLM(ctx context.Context) ([]model.BIHistory, error)
	CreateBIHistory(ctx context.Context, bIHistory model.BIHistory) (int, error)
	UpdateBIHistory(ctx context.Context, bIHistory model.BIHistory) (int, error)
}

type Storage struct {
	IUserStorage
	IBookStorage
	IBIHistoryStorage
}

func NewStorage(logger *zap.Logger, db *sqlx.DB) *Storage {
	return &Storage{
		IUserStorage:      postgres.NewUserStorage(db, logger),
		IBookStorage:      postgres.NewBookStorage(db, logger),
		IBIHistoryStorage: postgres.NewBIHistory(db, logger),
	}
}
