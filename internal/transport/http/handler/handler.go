package handler

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/logger"
	"net/http"
	"time"
)

const _timeoutContext = 5 * time.Second

type IUser interface {
	GetUserByName(name string) (model.User, error)
	GetAllUsers() ([]*model.User, error)
	CreateUser(user model.User) (string, error)
	UpdateUser(name string, user model.User) error
	DeleteUser(name string) (model.User, error)
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

	name, err := h.user.CreateUser(*user)
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

	w.Write([]byte("User saved: " + name))
}

func (h *Handler) showUser(w http.ResponseWriter, r *http.Request) {
	name := httprouter.ParamsFromContext(r.Context()).ByName("name")

	user, err := h.user.GetUserByName(name)
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

	jsons := make([]byte, 0, len(users))

	for _, user := range users {
		json, err := h.formatToJSON(*user)
		if err != nil {
			h.serverError(w, err)
			return
		}

		jsons = append(jsons, json...)
		jsons = append(jsons, '\n')
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsons)
}

func (h *Handler) updateCreateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	name := httprouter.ParamsFromContext(r.Context()).ByName("name")

	user, err := h.getUserDate(r)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err = h.user.UpdateUser(name, *user); err != nil {
		if errors.Is(err, common.ErrInvalidData) {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		name, _ = h.user.CreateUser(*user)

	}

	w.Write([]byte(name + " updated"))
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	name := httprouter.ParamsFromContext(r.Context()).ByName("name")

	user, err := h.getUserDate(r)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err = h.user.UpdateUser(name, *user); err != nil {
		if errors.Is(err, common.ErrUserNotExists) {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		h.serverError(w, err)
		return
	}

	w.Write([]byte(name + "user updated"))
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	name := httprouter.ParamsFromContext(r.Context()).ByName("name")

	user, err := h.user.DeleteUser(name)
	if err != nil {
		if errors.Is(err, common.ErrUserNotExists) {
			h.notFound(w)
			return
		}

		h.serverError(w, err)
	}

	json, err := h.formatToJSON(user)
	if err != nil {
		h.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
