package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/zhayt/transaction-service/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

const _timeoutContext = 5 * time.Second

// CreateTransaction godoc
// @Summary		Create transaction
// @Tags			transaction
// @Description	create transaction
// @ID				create-transaction
// @Accept			json
// @Param			input	body	model.Transaction	true	"transaction date"
// @Success		200		{int} transactionID
// @Failure		400		""
// @Failure		500		""
// @Router			/transactions [post]
func (h *Handler) CreateTransaction(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	var transaction model.Transaction

	if err := e.Bind(&transaction); err != nil {
		h.l.Error("Bind error", zap.Error(err))
		return e.NoContent(http.StatusBadRequest)
	}

	transactionID, err := h.transaction.CreateTransaction(ctx, transaction)
	if err != nil {
		h.l.Error("Create transaction error", zap.Error(err))
		return e.NoContent(http.StatusInternalServerError)
	}

	h.l.Info("Transaction has been created", zap.Int("id", transactionID))
	return e.JSON(http.StatusOK, transactionID)
}

// CreateTransactionItem godoc
// @Summary		Create transaction item
// @Tags			transaction
// @Description	create transaction items
// @ID				create-transactionItem
// @Accept			json
// @Param			input	body	model.TransactionItem	true	"transaction item date"
// @Success		200		""
// @Failure		400		""
// @Failure		500		""
// @Router			/transactions/items [post]
func (h *Handler) CreateTransactionItem(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	var item model.TransactionItem

	if err := e.Bind(&item); err != nil {
		h.l.Error("Bind error", zap.Error(err))
		return e.NoContent(http.StatusBadRequest)
	}

	itemID, err := h.transaction.CreateTransactionItem(ctx, item)
	if err != nil {
		h.l.Error("Create transaction item error", zap.Error(err))
		return e.NoContent(http.StatusInternalServerError)
	}

	h.l.Info("Transaction item has been created", zap.Int("itemID", itemID), zap.Uint("transactionID", item.TransactionID))
	return e.NoContent(http.StatusOK)
}

// DeleteTransaction godoc
// @Summary		delete transaction
// @Tags			transaction
// @Description	delete transaction with items from db
// @ID				delete-transaction
// @Param			id	path	integer	true	"transactionID"
// @Success		200	""
// @Failure		404	""
// @Failure		500	""
// @Router			/transaction/{id} [delete]
func (h *Handler) DeleteTransaction(e echo.Context) error {
	ctx, cancel := context.WithTimeout(e.Request().Context(), _timeoutContext)
	defer cancel()

	transactionID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		h.l.Error("Param error", zap.Error(err))
		return e.NoContent(http.StatusNotFound)
	}

	if err = h.transaction.DeleteTransaction(ctx, transactionID); err != nil {
		h.l.Error("Delete transaction error", zap.Error(err))
		return e.NoContent(http.StatusInternalServerError)
	}

	h.l.Info("Transaction has been deleted", zap.Int("id", transactionID))
	return e.NoContent(http.StatusOK)
}
