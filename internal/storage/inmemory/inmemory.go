package inmemory

import (
	"fmt"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"sync"
)

type UserStorage struct {
	mu      sync.Mutex
	storage map[string]model.User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		storage: make(map[string]model.User, 50),
	}
}

func (r *UserStorage) GetUserByName(name string) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.storage[name]
	if !ok {
		return user, fmt.Errorf("cannot take user: %w", common.ErrUserNotExists)
	}

	return user, nil
}

func (r *UserStorage) GetAllUsers() ([]*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	users := make([]*model.User, 0, len(r.storage))
	for _, user := range r.storage {
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserStorage) CreateUser(name string, user model.User) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.storage[name]; ok {
		return "", fmt.Errorf("cannot create user: %w", common.ErrNameTaken)
	}

	r.storage[name] = user

	return name, nil
}

func (r *UserStorage) UpdateUser(name string, user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[name]; !ok {
		return common.ErrUserNotExists
	}

	r.storage[name] = user

	return nil
}

func (r *UserStorage) DeleteUser(name string) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.storage[name]
	if !ok {
		return user, common.ErrUserNotExists
	}

	delete(r.storage, name)
	return user, nil
}
