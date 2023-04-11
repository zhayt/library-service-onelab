package handler

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
	"net/http"
)

type IUserService interface {
	GetUserById(ctx context.Context, userID int) (model.User, error)
	GetUserByEmail(ctx context.Context, login model.UserLogin) (model.User, error)

	CreateUser(ctx context.Context, user model.User) (model.User, error)

	UpdateUserFIO(ctx context.Context, user model.UserUpdateFIO) (int, error)
	UpdateUserPassword(ctx context.Context, user model.UserUpdatePassword) (int, error)

	DeleteUser(ctx context.Context, userID int) error
}

func (h *Handler) CreateUser(e echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), _timeoutContext)
	defer cancel()

	var user model.User
	if err := e.Bind(&user); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	user, err := h.user.CreateUser(ctx, user)
	if err != nil {
		h.log.Error("Create error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	h.log.Info("User created", zap.Int("id", user.ID))

	return e.JSON(http.StatusOK, user)
}

func (h *Handler) GetUser(e echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), _timeoutContext)
	defer cancel()

	var userLogin model.UserLogin
	if err := e.Bind(&userLogin); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	user, err := h.user.GetUserByEmail(ctx, userLogin)
	if err != nil {
		log.Error("Get user error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	token, err := h.mid.GenerateJWT(user.FIO, user.ID)
	if err != nil {
		h.log.Error("Generate token error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
func (h *Handler) UpdateUser(e echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), _timeoutContext)
	defer cancel()

	userId, err := getUserId(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, err)
	}

	var userdata model.UserUpdateFIO

	if err = e.Bind(&userdata); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	userdata.ID = userId

	userId, err = h.user.UpdateUserFIO(ctx, userdata)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	return e.JSON(http.StatusOK, userId)
}

func (h *Handler) UpdateUserPassword(e echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), _timeoutContext)
	defer cancel()

	userId, err := getUserId(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, err)
	}

	var userPasswd model.UserUpdatePassword

	if err = e.Bind(&userPasswd); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	userPasswd.ID = userId

	userId, err = h.user.UpdateUserPassword(ctx, userPasswd)
	if err != nil {
		// will be Server or Client error
		return e.JSON(http.StatusBadRequest, err)
	}

	return e.JSON(http.StatusOK, userId)
}

func (h *Handler) DeleteUser(e echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), _timeoutContext)
	defer cancel()

	userId, err := getUserId(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, err)
	}

	if err = h.user.DeleteUser(ctx, userId); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	return e.JSON(http.StatusOK, userId)
}

func getUserId(e echo.Context) (int, error) {
	id := e.Get(model.ContextUserID)

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
