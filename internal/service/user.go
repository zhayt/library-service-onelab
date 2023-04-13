package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"strings"
)

type IUserStorage interface {
	GetUserByID(ctx context.Context, userID int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)

	CreateUser(ctx context.Context, user model.User) (int, error)

	UpdateUserFIO(ctx context.Context, user model.UserUpdateFIO) (int, error)
	UpdateUserPassword(ctx context.Context, user model.UserUpdatePassword) (int, error)

	DeleteUser(ctx context.Context, userID int) error
}

type UserService struct {
	user storage.IUserStorage
	log  *zap.Logger
}

func NewUserService(logger *zap.Logger, user IUserStorage) *UserService {
	return &UserService{log: logger, user: user}
}

func (s *UserService) GetUserByID(ctx context.Context, userId int) (model.User, error) {
	return s.user.GetUserByID(ctx, userId)
}

func (s *UserService) GetUserByEmail(ctx context.Context, login model.UserLogin) (model.User, error) {
	user, err := s.user.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return user, err
	}

	if err = compareHashAndPassword(user.Password, login.Password); err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) (int, error) {
	if err := checkDate(3, 50, user.FIO, user.Password); err != nil {
		s.log.Error("Check date error", zap.Error(err))
		return 0, err
	}

	if err := matchesPattern(user.Email, common.EmailRX); err != nil {
		return 0, err
	}

	passwdHash, err := generatePasswordHash(user.Password)
	if err != nil {
		s.log.Error("Passwd hash error", zap.Error(err))
		return 0, err
	}

	user.Password = passwdHash

	userID, err := s.user.CreateUser(ctx, user)
	if err != nil {
		s.log.Error("Service error", zap.Error(err))
		return 0, err
	}

	return userID, nil
}

func (s *UserService) UpdateUserFIO(ctx context.Context, user model.UserUpdateFIO) (int, error) {
	if err := checkDate(3, 50, user.FIO); err != nil {
		return 0, err
	}

	user.FIO = strings.TrimSpace(user.FIO)

	return s.user.UpdateUserFIO(ctx, user)
}

func (s *UserService) UpdateUserPassword(ctx context.Context, userUP model.UserUpdatePassword) (int, error) {
	if err := checkDate(3, 50, userUP.NewPassword); err != nil {
		return 0, err
	}

	// Насколько это правельно вызывать внутни сервиса другой сервис???
	user, err := s.GetUserByID(ctx, userUP.ID)
	if err != nil {
		return 0, err
	}

	if err = compareHashAndPassword(user.Password, userUP.CurrentPassword); err != nil {
		return 0, err
	}

	if userUP.NewPassword != userUP.NewPasswordRepeat {
		return 0, errors.New("not same password")
	}

	passwdHash, err := generatePasswordHash(userUP.NewPassword)
	if err != nil {
		return 0, err
	}

	userUP.NewPassword = passwdHash

	return s.user.UpdateUserPassword(ctx, userUP)
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	// Можно сделать чтобы пользователь ввел пароль
	// и проверять сответствие пароля перед тем удалять пользователя
	return s.user.DeleteUser(ctx, userID)
}

func generatePasswordHash(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't generate password hash: %w", err)
	}

	return string(hash), nil
}

func compareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
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
