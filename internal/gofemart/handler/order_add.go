package handler

import (
	"context"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order/status"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
)

func (h *APIhandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	logger := logging.Logger.GetSubLogger("Method", "Order Add")
	logger.Warn("Processing request")
	ctx := r.Context()
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
	logger = logger.GetSubLogger("orderID", orderID)
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

	err = h.storage.WithinTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err = h.storage.GetOrder(ctx, orderID, tx)
		switch {
		case errors.Is(err, repositoryerrors.ErrNoContent):
		case err != nil:
			logger.Warn("error during fetching order")
			http.Error(w, "System error", http.StatusInternalServerError)
			return err
		default:
			owner, err := h.storage.GetOrderOwner(ctx, orderID, tx)
			switch {
			case owner == userID:
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				return err
			case err != nil:
				logger.Warn("error during fetching order")
				http.Error(w, "System error", http.StatusInternalServerError)
				return err
			default:
				logger.Warn("Order was already registered in order database")
				http.Error(w, "Order was already registered", http.StatusConflict)
				return err
			}
		}
		_, err = h.storage.GetWithdrawn(ctx, orderID, tx)
		switch {
		case errors.Is(err, repositoryerrors.ErrNoContent):
		case err != nil:
			logger.Warn("error during fetching Withdraw: %s", err.Error())
			http.Error(w, "Order %s already exists", http.StatusInternalServerError)
			return err
		default:
			logger.Warn("Order was already registered in withdraw database")
			http.Error(w, "Order was already registered", http.StatusConflict)
			return err
		}
		logger.Warn("Trying to bind order=%s to user-account=%s", orderID, userID)
		accuralRecord, err := h.accrualClient.FetchOrderInfo(ctx, orderID)
		if err != nil {
			logger.Warn("Failed to fecth accural Info: %s", err.Error())
			http.Error(w, "error during registering order", http.StatusInternalServerError)
			return err
		}

		order := order.FromAccural(accuralRecord)

		err = h.storage.AddOrder(ctx, userID, order, tx)
		if err != nil {
			logger.Warn("error during adding order to order Database: %s", err.Error())
			http.Error(w, "error during registering order", http.StatusInternalServerError)
			return err
		}
		if status.RequiresTracking(order.Status) {
			logger.Debug("adding order to PENDING")
			err = h.storage.AddOrderToPendingList(ctx, orderID, tx)
			if err != nil {
				logger.Warn("error during adding order to pending list: %s", err.Error())
				http.Error(w, "error during registering order", http.StatusInternalServerError)
				return err
			}
		}
		if accuralRecord.Status == status.PROCESSED {
			logger.Debug("Applying accural to balance")
			accountData, err := h.storage.GetCustomerAccount(ctx, userID, tx)
			if err != nil {
				logger.Warn("error during fetching account info: %s", err.Error())
				http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
				return err
			}
			accountData.Current += order.Accural
			err = h.storage.UpdateCustomerAccount(ctx, accountData, tx)
			if err != nil {
				logger.Warn("error during fetching account info: %s", err.Error())
				http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Debug("New Order registration failed")
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	logger.Debug("New Order is successfully registered")
}
