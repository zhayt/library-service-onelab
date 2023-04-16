package storage

import (
	"context"
	"github.com/zhayt/user-storage-service/config"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/storage/postgres"
	"go.uber.org/zap"
	"sync"
)

type IUserStorage interface {
	GetUserByID(ctx context.Context, userID int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (int, error)
	UpdateUserFIO(ctx context.Context, user model.UserUpdateFIO) (int, error)
	UpdateUserPassword(ctx context.Context, user model.UserUpdatePassword) (int, error)
	DeleteUser(ctx context.Context, userID int) error
}

type IBookStorage interface {
	GetBookByID(ctx context.Context, bookID int) (model.Book, error)
	GetAllBooks(ctx context.Context) ([]model.Book, error)
	CreateBook(ctx context.Context, book model.Book) (int, error)
	UpdateBook(ctx context.Context, book model.Book) (int, error)
	DeleteBook(ctx context.Context, bookID int) error
}

type IBIHistoryStorage interface {
	GetCurrentBorrowedBooks(ctx context.Context) ([]model.BorrowedBooks, error)
	GetBIHistoryLastMonth(ctx context.Context) ([]model.BorrowedBooks, error)
	CreateBIHistory(ctx context.Context, bIHistory model.BIHistory) (int, error)
	UpdateBIHistory(ctx context.Context, bIHistoryID int) (int, error)
	DeleteBIHistory(ctx context.Context, bIHistoryID int) error
}

type Storage struct {
	IUserStorage
	IBookStorage
	IBIHistoryStorage
}

func NewStorage(ctx context.Context, wg *sync.WaitGroup, logger *zap.Logger, cfg *config.Config) (*Storage, error) {
	db, err := postgres.Dial("pgx", cfg.DBConnectionURL)
	if err != nil {
		logger.Error("Dial error", zap.Error(err))
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			logger.Info("Close db pool connection")
			if err := db.Close(); err != nil {
				logger.Error("Close pool connection error", zap.Error(err))
			}
		}
	}()

	return &Storage{
		IUserStorage:      postgres.NewUserStorage(db, logger),
		IBookStorage:      postgres.NewBookStorage(db, logger),
		IBIHistoryStorage: postgres.NewBIHistory(db, logger),
	}, nil
}
