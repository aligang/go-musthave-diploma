package handler

import (
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/accural"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order/status"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *APIhandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	logging.Warn("Processing order add request")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	userID, err := auth.ResolveUsername(r)
	logging.Debug("Processing Order registration request for user %s", userID)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logging.Warn("No user info were provided")
		return
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read data from wire", http.StatusInternalServerError)
		return
	}
	orderID := string(payload)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = order.ValidateID(orderID)
	if err != nil {
		logging.Warn("Invalid order format: %s", orderID)
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}
	err = order.ValidateIDFormat(orderID)
	if err != nil {
		logging.Warn("Invalid orderID checksum: %s", orderID)
		http.Error(w, "Invalid orderID checksum", http.StatusUnprocessableEntity)
		return
	}
	logging.Debug("Processing Order registration request with id %s", orderID)
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
		logging.Warn("error during fetching order with id: %s", orderID)
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
			logging.Warn("error during fetching order=%s", orderID)
			http.Error(w, "System error", http.StatusInternalServerError)
			return
		default:
			logging.Warn("Order=%s was already registered in order database", orderID)
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
		logging.Warn("error during fetching Withdraw %s", err.Error())
		http.Error(w, "Account %s already exists", http.StatusInternalServerError)
		return
	default:
		logging.Warn("Order=%s was already registered in withdraw database", orderID)
		http.Error(w, "Order was already registered", http.StatusConflict)
		return
	}

	logging.Warn("Trying to bind order=%s to user-account=%s", orderID, userID)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	accuralRecord, err := accural.FetchOrderInfo(ctx, orderID, h.config)
	if err != nil {
		logging.Warn("Failed to fecth accural Info: %s", err.Error())
		http.Error(w, "error during registering order", http.StatusInternalServerError)
		return
	}

	order := order.FromAccural(accuralRecord)

	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = h.storage.AddOrder(ctx, userID, order)
	if err != nil {
		logging.Warn("error during adding order to order Database: %s", err.Error())
		http.Error(w, "error during registering order", http.StatusInternalServerError)
		return
	}
	if status.RequiresTracking(order.Status) {
		logging.Debug("adding order=%s to PENDING", order.Number)
		if RequestContextIsClosed(ctx, w) {
			return
		}
		err = h.storage.AddOrderToPendingList(ctx, orderID)
		if err != nil {
			logging.Warn("error during adding order to pending list: %s", err.Error())
			http.Error(w, "error during registering order", http.StatusInternalServerError)
			return
		}
	}
	if accuralRecord.Status == status.PROCESSED {
		logging.Debug("Applying orderID=%s accural to %s balance", accuralRecord.Order, userID)
		if RequestContextIsClosed(ctx, w) {
			return
		}
		accountData, err := h.storage.GetCustomerAccountWithinTransaction(ctx, userID)
		if err != nil {
			logging.Warn("error during fetching account info: %s", err.Error())
			http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
			return
		}
		accountData.Current += order.Accural
		if RequestContextIsClosed(ctx, w) {
			return
		}
		err = h.storage.UpdateCustomerAccount(ctx, accountData)
		if err != nil {
			logging.Warn("error during fetching account info: %s", err.Error())
			http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
			return
		}
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	logging.Debug("New Order=%s is successfully registered", orderID)
}
