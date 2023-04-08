package service

import (
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"net/mail"
	"regexp"
	"strings"
)

type IUserService interface {
	GetUserById(id int) (model.User, error)
	GetAllUsers() ([]model.User, error)
	CreateUser(user model.User) (int, error)
	UpdateUser(id int, user model.User) error
	DeleteUser(id int) error
}

type UserService struct {
	User IUserService
}

// Ты принимаешь и отдаешь одну и ту же структуру / интерфейс ? тут правильнее будет использовать IUserStorage
func NewUserService(user IUserService) *UserService {
	return &UserService{User: user}
}

func (s *UserService) GetUserById(id int) (model.User, error) {
	return s.User.GetUserById(id)
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.User.GetAllUsers()
}

func (s *UserService) CreateUser(user model.User) (int, error) {
	if err := checkDate(3, 50, user.FIO, user.PasswordHash); err != nil {
		return 0, err
	}

	if err := matchesPattern(user.Email, common.EmailRX); err != nil {
		return 0, err
	}

	// generate password hash

	user.FIO = strings.TrimSpace(user.FIO)
	user.PasswordHash = strings.TrimSpace(user.PasswordHash)

	return s.User.CreateUser(user)
}

func (s *UserService) UpdateUser(id int, user model.User) error {
	if err := checkDate(3, 50, user.FIO, user.PasswordHash); err != nil {
		return err
	}

	if err := matchesPattern(user.Email, common.EmailRX); err != nil {
		return err
	}

	// generate password hash

	user.FIO = strings.TrimSpace(user.FIO)
	user.PasswordHash = strings.TrimSpace(user.PasswordHash)

	return s.User.UpdateUser(id, user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.User.DeleteUser(id)
}

func checkDate(minLen, maxLen int, data ...string) error {
	for _, d := range data {
		d = strings.TrimSpace(d)
		if len([]rune(d)) < minLen && len([]rune(d)) > maxLen {
			return common.ErrInvalidData
		}
	}

	return nil
}

func matchesPattern(value string, pattern *regexp.Regexp) error {
	if value == "" {
		return common.ErrInvalidData
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		return common.ErrInvalidData
	}
	if !pattern.MatchString(value) {
		return common.ErrInvalidData
	}

	return nil
}
