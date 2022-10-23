package handler

import (
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order/status"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *APIhandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	logger := logging.Logger.GetSubLogger("Method", "Order Add")
	logger.Warn("Processing request")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	userID, err := auth.ResolveUsername(r)
	logger = logger.GetSubLogger("userID", userID)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logger.Warn("No user info were provided")
		return
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read data from wire", http.StatusInternalServerError)
		return
	}
	orderID := string(payload)
	logger = logger.GetSubLogger("orderId", orderID)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = order.ValidateID(orderID)
	if err != nil {
		logger.Warn("Invalid order format")
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}
	err = order.ValidateIDFormat(orderID)
	if err != nil {
		logger.Warn("Invalid orderID checksum")
		http.Error(w, "Invalid orderID checksum", http.StatusUnprocessableEntity)
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	var dBerr error
	h.storage.StartTransaction(ctx)
	defer func() {
		if dBerr != nil {
			h.storage.RollbackTransaction(ctx)
		}
		h.storage.CommitTransaction(ctx)
	}()
	_, err = h.storage.GetOrderWithinTransaction(ctx, orderID)
	switch {
	case errors.Is(err, repositoryerrors.ErrNoContent):
	case err != nil:
		logger.Warn("error during fetching order")
		http.Error(w, "System error", http.StatusInternalServerError)
		return
	default:
		owner, err := h.storage.GetOrderOwner(ctx, orderID)
		switch {
		case owner == userID:
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			return
		case err != nil:
			logger.Warn("error during fetching order")
			http.Error(w, "System error", http.StatusInternalServerError)
			return
		default:
			logger.Warn("Order was already registered in order database")
			http.Error(w, "Order was already registered", http.StatusConflict)
			return
		}
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	_, err = h.storage.GetWithdrawnWithinTransaction(ctx, orderID)
	switch {
	case errors.Is(err, repositoryerrors.ErrNoContent):
	case err != nil:
		logger.Warn("error during fetching Withdraw: %s", err.Error())
		http.Error(w, "Order %s already exists", http.StatusInternalServerError)
		return
	default:
		logger.Warn("Order was already registered in withdraw database")
		http.Error(w, "Order was already registered", http.StatusConflict)
		return
	}

	logger.Warn("Trying to bind order=%s to user-account=%s", orderID, userID)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	accuralRecord, err := h.accrualClient.FetchOrderInfo(ctx, orderID)
	if err != nil {
		logger.Warn("Failed to fecth accural Info: %s", err.Error())
		http.Error(w, "error during registering order", http.StatusInternalServerError)
		return
	}

	order := order.FromAccural(accuralRecord)

	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = h.storage.AddOrder(ctx, userID, order)
	if err != nil {
		logger.Warn("error during adding order to order Database: %s", err.Error())
		http.Error(w, "error during registering order", http.StatusInternalServerError)
		return
	}
	if status.RequiresTracking(order.Status) {
		logger.Debug("adding order to PENDING")
		if RequestContextIsClosed(ctx, w) {
			return
		}
		err = h.storage.AddOrderToPendingList(ctx, orderID)
		if err != nil {
			logger.Warn("error during adding order to pending list: %s", err.Error())
			http.Error(w, "error during registering order", http.StatusInternalServerError)
			return
		}
	}
	if accuralRecord.Status == status.PROCESSED {
		logger.Debug("Applying accural to balance")
		if RequestContextIsClosed(ctx, w) {
			return
		}
		accountData, err := h.storage.GetCustomerAccountWithinTransaction(ctx, userID)
		if err != nil {
			logger.Warn("error during fetching account info: %s", err.Error())
			http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
			return
		}
		accountData.Current += order.Accural
		if RequestContextIsClosed(ctx, w) {
			return
		}
		err = h.storage.UpdateCustomerAccount(ctx, accountData)
		if err != nil {
			logger.Warn("error during fetching account info: %s", err.Error())
			http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
			return
		}
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	logger.Debug("New Order is successfully registered")
}
