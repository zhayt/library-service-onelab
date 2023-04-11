package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/zhayt/user-storage-service/internal/model"
	"net/http"
)

type IBookService interface {
	CreateBook(ctx context.Context, book model.Book) (model.Book, error)
	GetBookById(ctx context.Context, bookId int) (model.Book, error)
	GetBooks(ctx context.Context) ([]model.Book, error)
	UpdateBook(ctx context.Context, book model.Book) (int, error)
	DeleteBook(ctx context.Context, bookId int) error
}

func (h *Handler) CreateBook(e echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), _timeoutContext)
	defer cancel()

	var book model.Book

	if err := e.Bind(&book); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	book, err := h.book.CreateBook(ctx, book)
	if err != nil {
		// server or client error
		return e.JSON(http.StatusBadRequest, err)
	}

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
