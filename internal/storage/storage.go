package storage

import (
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/storage/inmemory"
)

type IUserStorage interface {
	GetUserById(id int) (model.User, error)
	GetAllUsers() ([]model.User, error)
	CreateUser(user model.User) (int, error)
	UpdateUser(id int, user model.User) error
	DeleteUser(id int) error
}

type Storage struct {
	IUserStorage
}

func NewStorage() *Storage {
	storage := inmemory.NewUserStorage()
	return &Storage{storage}
}
