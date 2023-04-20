package handler

import (
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/service"
	"github.com/zhayt/user-storage-service/internal/transport/http/middleware"
	"go.uber.org/zap"
	"time"
)

const _timeoutContext = 5 * time.Second

type Handler struct {
	log         *zap.Logger
	user        IUserService
	book        IBookService
	history     IBIHistoryService
	rent        service.IRentTransactionService
	transaction service.ITransactionService
	mid         *middleware.JWTAuth
}

func NewHandler(logger *zap.Logger, service *service.Service, auth *middleware.JWTAuth) *Handler {
	return &Handler{
		log:     logger,
		user:    service,
		book:    service,
		rent:    service,
		history: service,
		mid:     auth,
	}
}

func makeResponse(msg interface{}) *model.Response {
	return &model.Response{Message: msg}
}
