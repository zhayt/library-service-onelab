package handler

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/logger"
	"net/http"
	"strconv"
	"time"
)

const _timeoutContext = 5 * time.Second

const (
	responseUserSave   = "User Saved"
	responseUserUpdate = "User Updated"
	responseUserDelete = "User Deleted"
)

type IUser interface {
	GetUserById(id int) (model.User, error)
	GetAllUsers() ([]model.User, error)
	CreateUser(user model.User) (int, error)
	UpdateUser(id int, user model.User) error
	DeleteUser(id int) error
}

type Handler struct {
	*logger.Logger
	user IUser
}

func NewHandler(logger *logger.Logger, user IUser) *Handler {
	return &Handler{
		logger,
		user}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := h.getUserDate(r)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := h.user.CreateUser(*user)
	if err != nil {
		if errors.Is(err, common.ErrInvalidData) {
			h.clientError(w, http.StatusBadRequest)
			return
		}
		if errors.Is(err, common.ErrNameTaken) {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		h.serverError(w, err)
		return
	}

	response := model.Response{
		Message: fmt.Sprintf("%v: #%v", responseUserSave, id),
	}

	json, err := h.formatToJSON(response)
	if err != nil {
		h.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (h *Handler) showUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := h.user.GetUserById(id)
	if err != nil {
		h.notFound(w)
		return
	}

	json, err := h.formatToJSON(user)
	if err != nil {
		h.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (h *Handler) showAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.user.GetAllUsers()
	if err != nil {
		h.serverError(w, err)
		return
	}

	json, err := h.formatToJSON(users)
	if err != nil {
		h.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (h *Handler) updateCreateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := h.getUserDate(r)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err = h.user.UpdateUser(id, *user); err != nil {
		if errors.Is(err, common.ErrInvalidData) {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		id, _ = h.user.CreateUser(*user)

	}

	response := model.Response{
		Message: fmt.Sprintf("%v: #%v", responseUserUpdate, id),
	}
	json, err := h.formatToJSON(response)
	if err != nil {
		h.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := h.getUserDate(r)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err = h.user.UpdateUser(id, *user); err != nil {
		if errors.Is(err, common.ErrUserNotExists) {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		h.serverError(w, err)
		return
	}

	response := model.Response{
		Message: fmt.Sprintf("%v: #%v", responseUserUpdate, id),
	}
	json, err := h.formatToJSON(response)
	if err != nil {
		h.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	err = h.user.DeleteUser(id)
	if err != nil {
		if errors.Is(err, common.ErrUserNotExists) {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		h.serverError(w, err)
	}

	response := model.Response{
		Message: fmt.Sprintf("%v: #%v", responseUserDelete, id),
	}
	json, err := h.formatToJSON(response)
	if err != nil {
		h.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
