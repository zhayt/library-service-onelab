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
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bookID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, err)
	}

	book, err := h.book.GetBookByID(ctx, bookID)
	if err != nil {
		h.log.Error("Get book Id error", zap.Error(err))
		// 500 or 404
		return e.JSON(http.StatusNotFound, err)
	}

	h.log.Info("Book found", zap.Int("id", book.ID))
	return e.JSON(http.StatusOK, book)
}

func (h *Handler) ShowAllBooks(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	books, err := h.book.GetAllBooks(ctx)
	if err != nil {
		h.log.Error("Get all books error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, err)
	}

func (h *Handler) ShowBooks(e echo.Context) error {
	return nil
}

func (h *Handler) UpdateBook(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bookID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, err)
	}

	var book model.Book
	if err := e.Bind(&book); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	book.ID = bookID

	bookID, err = h.book.UpdateBook(ctx, book)
	if err != nil {
		h.log.Error("Update book error", zap.Error(err))
		// 500 or 400
		return e.JSON(http.StatusInternalServerError, err)
	}

	h.log.Info("Book updated", zap.Int("id", bookID))
	return e.JSON(http.StatusOK, bookID)
}

func (h *Handler) DeleteBook(e echo.Context) error {
	return nil
}
