package inmemory

import (
	"context"
	"fmt"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"sync"
)

type UserStorage struct {
	mu      sync.RWMutex
	id      int
	storage map[int]model.User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		storage: make(map[int]model.User, 50),
	}
}

func (r *UserStorage) GetUserById(ctx context.Context, id int) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.storage[id]
	if !ok {
		return user, fmt.Errorf("cannot take user: %w", common.ErrUserNotExists)
	}

	return user, nil
}

func (r *UserStorage) GetAllUsers(ctx context.Context) ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]model.User, 0, len(r.storage))
	for _, user := range r.storage {
		users = append(users, user)
	}

	return users, nil
}

func (r *UserStorage) CreateUser(ctx context.Context, user model.User) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id++
	if _, ok := r.storage[r.id]; ok {
		return 0, fmt.Errorf("cannot create user: %w", common.ErrNameTaken)
	}

	r.storage[r.id] = user

	return r.id, nil
}

func (r *UserStorage) UpdateUser(ctx context.Context, user model.User) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[user.ID]; !ok {
		return 0, fmt.Errorf("cannot update user: %w", common.ErrUserNotExists)
	}

	r.storage[user.ID] = user

	return user.ID, nil
}

func (r *UserStorage) DeleteUser(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[id]; !ok {
		return fmt.Errorf("cannot delete user: %w", common.ErrUserNotExists)
	}

	delete(r.storage, id)
	return nil
}
