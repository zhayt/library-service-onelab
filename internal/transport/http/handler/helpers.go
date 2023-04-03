package handler

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhayt/user-storage-service/internal/common"
	"github.com/zhayt/user-storage-service/internal/model"
	"log"
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
	firstname, ok := r.Form["firstname"]
	if !ok {
		log.Println("1")
		return nil, common.ErrEmptyField
	}

	lastname, ok := r.Form["lastname"]
	if !ok {
		log.Println("2")
		return nil, common.ErrEmptyField
	}

	fatherland, ok := r.Form["fatherland"]
	if !ok {
		log.Println("3")
		return nil, common.ErrEmptyField
	}

	email, ok := r.Form["email"]
	if !ok {
		log.Println("4")
		return nil, common.ErrEmptyField
	}

	password, ok := r.Form["password"]
	if !ok {
		log.Println("5")
		return nil, common.ErrEmptyField
	}

	user := &model.User{
		FirstName:    firstname[0],
		LastName:     lastname[0],
		Fatherland:   fatherland[0],
		Email:        email[0],
		PasswordHash: password[0],
	}

	return user, nil
}

func (h *Handler) formatToJSON(user model.User) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	js, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("cannot format to json: %w", err)
	}

	return js, nil
}
