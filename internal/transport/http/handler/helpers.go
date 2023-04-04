package handler

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"net/http"
	"runtime/debug"
)

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.LogError.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Handler) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}

func (h *Handler) getUserDate(r *http.Request) (*model.User, error) {
	fio, ok := r.Form["fio"]
	if !ok {
		return nil, fmt.Errorf("%w: field fio", common.ErrEmptyField)
	}

	email, ok := r.Form["email"]
	if !ok {
		return nil, fmt.Errorf("%w: field email", common.ErrEmptyField)
	}

	password, ok := r.Form["password"]
	if !ok {
		return nil, fmt.Errorf("%w: field password", common.ErrEmptyField)
	}

	user := &model.User{
		FIO:          fio[0],
		Email:        email[0],
		PasswordHash: password[0],
	}

	return user, nil
}

func (h *Handler) formatToJSON(user interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	js, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("cannot format to json: %w", err)
	}

	return js, nil
}
