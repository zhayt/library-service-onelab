package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type IBIHistoryService interface {
	GetCurrentBorrowedBooks(ctx context.Context) ([]model.BorrowedBooks, error)
	GetBIHistoryLastMonth(ctx context.Context) ([]model.BorrowedBooks, error)
	UpdateBIHistory(ctx context.Context, bIHistoryID int) (int, error)
	DeleteBIHistory(ctx context.Context, bIHistoryID int) error
}

// CreateBIHistory godoc
// @Summary		rent book
// @Security	ApiKeyAuth
// @Tags		book-issue-history
// @Description	create note about rent book
// @ID			rent-book
// @Accept		json
// @Produce		json
// @Param		input	body		model.BIHistory	true	"book issue info"
// @Success		200		""
// @Success		401		{object}	model.Response
// @Failure		400		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router		/rents [post]
func (h *Handler) CreateBIHistory(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	userID, err := getUserID(e)
	if err != nil {
		h.log.Error("Authorization error", zap.Error(err))
		return e.JSON(http.StatusUnauthorized, makeResponse(err.Error()))
	}

	var bIHistory model.BIHistory
	if err = e.Bind(&bIHistory); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, makeResponse(err.Error()))
	}

	bIHistory.UserID = userID

	if err = h.rent.RentBook(ctx, bIHistory); err != nil {
		h.log.Error("Create book issue history error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Book issue history has been created")

	return e.NoContent(http.StatusOK)
}

// ShowCurrentBorrowedBooks godoc
// @Summary		show current borrowed books
// @Tags		book-issue-history
// @Description	show current borrowed books
// @ID			show-rent-book
// @Produce		json
// @Success		200		{object}	[]model.BorrowedBooks
// @Failure		500		{object}	model.Response
// @Router		/rents [get]
func (h *Handler) ShowCurrentBorrowedBooks(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	borrowedBooks, err := h.history.GetCurrentBorrowedBooks(ctx)
	if err != nil {
		h.log.Error("Get current borrowed books error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Showed current borrowed books", zap.Int("amount", len(borrowedBooks)))
	return e.JSON(http.StatusOK, borrowedBooks)
}

// ShowBIHistoryLastMonth godoc
// @Summary		show borrowed books in last month
// @Tags		book-issue-history
// @Description	show borrowed books in last month
// @ID			show-rent-book-lm
// @Produce		json
// @Success		200		{object}	[]model.BorrowedBooks
// @Failure		500		{object}	model.Response
// @Router		/rents/months [get]
func (h *Handler) ShowBIHistoryLastMonth(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	borrowedBooks, err := h.history.GetBIHistoryLastMonth(ctx)
	if err != nil {
		h.log.Error("Get book issue history last month error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Showed borrowed books in last month", zap.Int("amount", len(borrowedBooks)))
	return e.JSON(http.StatusOK, borrowedBooks)
}

// UpdateBIHistory godoc
// @Summary		update book issue history
// @Tags		book-issue-history
// @Description	update book issue history book returned
// @ID			update-biHistory
// @Produce		json
// @Param 		id	path		integer	true	"BIHistoryID"
// @Success		200		{object}	[]model.BorrowedBooks
// @Failure		404		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router		/rents/{id} [patch]
func (h *Handler) UpdateBIHistory(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bIHistoryID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	if _, err = h.history.UpdateBIHistory(ctx, bIHistoryID); err != nil {
		h.log.Error("Update book issue history error", zap.Int("id", bIHistoryID))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Book issue history has been updated", zap.Int("id", bIHistoryID))
	return e.JSON(http.StatusOK, makeResponse(bIHistoryID))
}

// DeleteBIHistory godoc
// @Summary		delete book issue history
// @Tags		book-issue-history
// @Description	delete book issue history
// @ID			delete-biHistory
// @Produce		json
// @Param		id	path		integer	true	"BIHistoryID"
// @Success		200		{object}	model.Response
// @Failure		404		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router		/rents/{id} [delete]
func (h *Handler) DeleteBIHistory(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bIHistoryID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, makeResponse(err.Error()))
	}

	if err = h.history.DeleteBIHistory(ctx, bIHistoryID); err != nil {
		h.log.Error("Delete book issue history error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, makeResponse(err.Error()))
	}

	h.log.Info("Book issue history has been deleted", zap.Int("id", bIHistoryID))
	return e.JSON(http.StatusOK, makeResponse(bIHistoryID))
}
