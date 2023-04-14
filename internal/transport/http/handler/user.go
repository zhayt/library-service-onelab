package handler

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type IUserService interface {
	GetUserByID(ctx context.Context, userID int) (model.User, error)
	GetUserByEmail(ctx context.Context, login model.UserLogin) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (int, error)
	UpdateUserFIO(ctx context.Context, user model.UserUpdateFIO) (int, error)
	UpdateUserPassword(ctx context.Context, user model.UserUpdatePassword) (int, error)
	DeleteUser(ctx context.Context, userID int) error
}

func (h *Handler) SignUp(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	var user model.User
	if err := e.Bind(&user); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	userID, err := h.user.CreateUser(ctx, user)
	if err != nil {
		h.log.Error("Create user error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	user.ID = userID

	h.log.Info("User created", zap.Int("id", userID))

	return e.JSON(http.StatusOK, user)
}

func (h *Handler) SignIn(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	var userLogin model.UserLogin
	if err := e.Bind(&userLogin); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	user, err := h.user.GetUserByEmail(ctx, userLogin)
	if err != nil {
		log.Error("Get user error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	token, err := h.mid.GenerateJWT(user.FIO, user.ID)
	if err != nil {
		h.log.Error("Generate token error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("User sign-in JWT created", zap.Int("id", user.ID))
	return e.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) ShowUser(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	userID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	user, err := h.user.GetUserByID(ctx, userID)
	if err != nil {
		h.log.Error("GetUserByID error", zap.Error(err))
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	h.log.Info("Show user", zap.Int("id", userID))
	return e.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUserFIO(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	userID, err := getUserID(e)
	if err != nil {
		h.log.Error("Authorization error", zap.Error(err))
		return e.JSON(http.StatusUnauthorized, makeResponse(err.Error()))
	}

	var userdata model.UserUpdateFIO

	if err = e.Bind(&userdata); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	userdata.ID = userID

	userID, err = h.user.UpdateUserFIO(ctx, userdata)
	if err != nil {
		h.log.Error("UpdateUserFIO error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	h.log.Info("User FIO has been changed", zap.Int("id", userID))
	return e.JSON(http.StatusOK, makeResponse(userID))
}

func (h *Handler) UpdateUserPassword(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	userID, err := getUserID(e)
	if err != nil {
		h.log.Error("Authorization error", zap.Error(err))
		return e.JSON(http.StatusUnauthorized, makeResponse(err.Error()))
	}

	var userPasswd model.UserUpdatePassword

	if err = e.Bind(&userPasswd); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	userPasswd.ID = userID

	userID, err = h.user.UpdateUserPassword(ctx, userPasswd)
	if err != nil {
		// will be Server or Client error
		h.log.Error("UpdateUserPassword error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	h.log.Info("User password has been changed", zap.Int("id", userID))
	return e.JSON(http.StatusOK, makeResponse(userID))
}

func (h *Handler) DeleteUser(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	userID, err := getUserID(e)
	if err != nil {
		h.log.Error("Authorization error", zap.Error(err))
		return e.JSON(http.StatusUnauthorized, makeResponse(err.Error()))
	}

	if err = h.user.DeleteUser(ctx, userID); err != nil {
		h.log.Error("DeleteUser error")
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	h.log.Info("User deleted", zap.Int("id", userID))
	return e.JSON(http.StatusOK, makeResponse(userID))
}

func getUserID(e echo.Context) (int, error) {
	id, ok := e.Request().Context().Value(model.ContextUserID).(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return id, nil
}
