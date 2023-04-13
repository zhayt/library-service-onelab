package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type IBookService interface {
	CreateBook(ctx context.Context, book model.Book) (int, error)
	GetBookByID(ctx context.Context, bookID int) (model.Book, error)
	GetAllBooks(ctx context.Context) ([]model.Book, error)
	UpdateBook(ctx context.Context, book model.Book) (int, error)
	DeleteBook(ctx context.Context, bookID int) error
}

func (h *Handler) CreateBook(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	var book model.Book

	if err := e.Bind(&book); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	bookID, err := h.book.CreateBook(ctx, book)
	if err != nil {
		// server or client error
		h.log.Error("Create book error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	book.ID = bookID

	h.log.Info("Book created", zap.Int("id", bookID))
	return e.JSON(http.StatusOK, book)
}

func (h *Handler) ShowBook(e echo.Context) error {
	return nil
}

func (h *Handler) ShowBooks(e echo.Context) error {
	return nil
}

func (h *Handler) UpdateBook(e echo.Context) error {
	return nil
}

func (h *Handler) DeleteBook(e echo.Context) error {
	return nil
}
