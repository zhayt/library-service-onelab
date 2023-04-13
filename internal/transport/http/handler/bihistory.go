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
	CreateBIHistory(ctx context.Context, history model.BIHistory) (int, error)
	GetCurrentBorrowedBooks(ctx context.Context) ([]model.BorrowedBooks, error)
	GetBIHistoryLastMonth(ctx context.Context) ([]model.BorrowedBooks, error)
	UpdateBIHistory(ctx context.Context, bIHistoryID int) (int, error)
	DeleteBIHistory(ctx context.Context, bIHistoryID int) error
}

func (h *Handler) CreateBIHistory(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	userID, err := getUserID(e)
	if err != nil {
		h.log.Error("Authorization error", zap.Error(err))
		return e.JSON(http.StatusUnauthorized, err)
	}

	var bIHistory model.BIHistory
	if err = e.Bind(&bIHistory); err != nil {
		h.log.Error("Bind error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	bIHistory.UserID = userID

	bIHistoryID, err := h.history.CreateBIHistory(ctx, bIHistory)
	if err != nil {
		h.log.Error("Create book issue history error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, err)
	}

	h.log.Info("Book issue history has been created", zap.Int("id", bIHistoryID))
	return e.JSON(http.StatusOK, bIHistoryID)
}

func (h *Handler) ShowCurrentBorrowedBooks(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	borrowedBooks, err := h.history.GetCurrentBorrowedBooks(ctx)
	if err != nil {
		h.log.Error("Get current borrowed books error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, err)
	}

	h.log.Info("Showed current borrowed books", zap.Int("amount", len(borrowedBooks)))
	return e.JSON(http.StatusOK, borrowedBooks)
}

func (h *Handler) ShowBIHistoryLastMonth(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	borrowedBooks, err := h.history.GetBIHistoryLastMonth(ctx)
	if err != nil {
		h.log.Error("Get book issue history last month error", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, err)
	}

	h.log.Info("Showed borrowed books in last month", zap.Int("amont", len(borrowedBooks)))
	return e.JSON(http.StatusOK, borrowedBooks)
}

func (h *Handler) UpdateBIHistory(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bIHistoryID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, err)
	}

	if _, err = h.history.UpdateBIHistory(ctx, bIHistoryID); err != nil {
		h.log.Error("Update book issue history error", zap.Int("id", bIHistoryID))
		return e.JSON(http.StatusInternalServerError, err)
	}

	h.log.Info("Book issue history has been updated", zap.Int("id", bIHistoryID))
	return e.JSON(http.StatusOK, bIHistoryID)
}

func (h *Handler) DeleteBIHistory(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	bIHistoryID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.log.Error("Param error", zap.Error(err))
		return e.JSON(http.StatusNotFound, err)
	}

	if err = h.history.DeleteBIHistory(ctx, bIHistoryID); err != nil {
		h.log.Error("Delete book issue history error", zap.Error(err))
		return e.JSON(http.StatusBadRequest, err)
	}

	h.log.Info("Book issue history has been deleted", zap.Int("id", bIHistoryID))
	return e.JSON(http.StatusOK, bIHistoryID)
}
