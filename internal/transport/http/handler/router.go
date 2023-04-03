package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) InitRoute() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/users/", h.showAllUsers)
	router.HandlerFunc(http.MethodGet, "/v1/users/:name", h.showUser)
	router.HandlerFunc(http.MethodPost, "/v1/users/", h.createUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/:name", h.updateCreateUser)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:name", h.updateUser)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:name", h.deleteUser)

	return h.logRequest(h.secureHeaders(router))
}
