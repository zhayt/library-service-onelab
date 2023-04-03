package storage

import (
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/storage/inmemory"
)

type IUserStorage interface {
	GetUserByName(name string) (model.User, error)
	GetAllUsers() ([]*model.User, error)
	CreateUser(name string, user model.User) (string, error)
	UpdateUser(name string, user model.User) error
	DeleteUser(name string) (model.User, error)
}

type Storage struct {
	IUserStorage
}

func NewStorage() *Storage {
	storage := inmemory.NewUserStorage()
	return &Storage{storage}
}
