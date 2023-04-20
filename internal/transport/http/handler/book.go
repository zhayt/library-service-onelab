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

// CreateBook godoc
// @Summary		Create-book
// @Tags		book
// @Description	create book
// @ID			create-book
// @Accept		json
// @Produce		json
// @Param		input	body		model.Book	true	"book info"
// @Success		200		{object}	model.Book
// @Failure		400		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router		/books [post]
func (h *Handler) CreateBook(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	var book model.Book

	if err := e.Bind(&book); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	bookID, err := h.book.CreateBook(ctx, book)
	if err != nil {
		// server or client error
		h.log.Error("Create book error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	book.ID = bookID

	h.log.Info("Book created", zap.Int("id", bookID))
	return e.JSON(http.StatusOK, book)
}

// ShowBook godoc
// @Summary		Show book
// @Tags		book
// @Description	show book
// @ID			show-book
// @Produce		json
// @Param		id	path		integer	true	"BookID"
// @Success		200		{object}	model.Book
// @Failure		400		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router		/books/{id} [get]
func (h *Handler) ShowBook(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bookID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	book, err := h.book.GetBookByID(ctx, bookID)
	if err != nil {
		h.log.Error("Get book Id error", zap.Error(err))
		// 500 or 404
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	h.log.Info("Book found", zap.Int("id", book.ID))
	return e.JSON(http.StatusOK, book)
}

// ShowAllBooks godoc
// @Summary		Show all books
// @Tags		book
// @Description	show books
// @ID			show-books
// @Produce		json
// @Success		200		{object}	[]model.Book
// @Failure		500		{object}	model.Response
// @Router		/books [get]
func (h *Handler) ShowAllBooks(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	books, err := h.book.GetAllBooks(ctx)
	if err != nil {
		h.log.Error("Get all books error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Books founded", zap.Int("amount", len(books)))
	return e.JSON(http.StatusOK, books)
}

// UpdateBook godoc
// @Summary		Update book
// @Tags		book
// @Description	update books
// @ID			update-book
// @Accept		json
// @Produce		json
// @Param		input	body		model.Book	true "book info"
// @Success		200		{object}	model.Book
// @Success		400		{object}	model.Book
// @Failure		404		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router		/books/{id} [patch]
func (h *Handler) UpdateBook(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bookID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	var book model.Book

	if err = e.Bind(&book); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	book.ID = bookID

	bookID, err = h.book.UpdateBook(ctx, book)
	if err != nil {
		h.log.Error("Update book error", zap.Error(err))
		// 500 or 400
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Book updated", zap.Int("id", bookID))
	return e.JSON(http.StatusOK, makeResponse(bookID))
}

// DeleteBook godoc
// @Summary		Delete book
// @Tags		book
// @Description	delete books
// @ID			delete-book
// @Produce		json
// @Param		id	path		integer	true	"BookID"
// @Success		200		{object}	model.Book
// @Success		404		{object}	model.Book
// @Failure		500		{object}	model.Response
// @Router		/books/{id} [delete]
func (h *Handler) DeleteBook(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bookID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	if err = h.book.DeleteBook(ctx, bookID); err != nil {
		h.log.Error("Delete book error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Book deleted", zap.Int("id", bookID))
	return e.JSON(http.StatusOK, makeResponse(bookID))
}
