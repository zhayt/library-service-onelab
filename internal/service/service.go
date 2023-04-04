package service

import (
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"net/mail"
	"regexp"
	"strings"
)

type IUserService interface {
	GetUserByName(name string) (model.User, error)
	GetAllUsers() ([]*model.User, error)
	CreateUser(user model.User) (string, error)
	UpdateUser(name string, user model.User) error
	DeleteUser(name string) (model.User, error)
}

type UserService struct {
	User IUserService
}

func NewUserService(user IUserService) *UserService {
	return &UserService{User: user}
}

func (s *UserService) GetUserByName(email string) (model.User, error) {
	return s.User.GetUserByName(email)
}

func (s *UserService) GetAllUsers() ([]*model.User, error) {
	return s.User.GetAllUsers()
}

func (s *UserService) CreateUser(user model.User) (string, error) {
	if err := checkDate(3, 50, user.FIO, user.PasswordHash); err != nil {
		return "", err
	}

	if err := matchesPattern(user.Email, common.EmailRX); err != nil {
		return "", err
	}

	// generate password hash

	user.FIO = strings.TrimSpace(user.FIO)
	user.PasswordHash = strings.TrimSpace(user.PasswordHash)

	return s.User.CreateUser(user)
}

func (s *UserService) UpdateUser(name string, user model.User) error {
	if err := checkDate(3, 50, user.FIO, user.PasswordHash); err != nil {
		return err
	}

	if err := matchesPattern(user.Email, common.EmailRX); err != nil {
		return err
	}

	// generate password hash

	user.FIO = strings.TrimSpace(user.FIO)
	user.PasswordHash = strings.TrimSpace(user.PasswordHash)

	return s.User.UpdateUser(name, user)
}

func (s *UserService) DeleteUser(name string) (model.User, error) {
	return s.User.DeleteUser(name)
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
