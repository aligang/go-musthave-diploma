package handler

import (
	"encoding/json"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"io"
	"net/http"
)

func (h *APIhandler) AddWithdraw(w http.ResponseWriter, r *http.Request) {
	logger := logging.Logger.GetSubLogger("Method", "Withdrawn New")
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
	logger.Warn("Received request %s", string(payload))
	if err != nil {
		http.Error(w, "Could not read data from wire", http.StatusInternalServerError)
		return
	}

	withdrawRequest := withdrawn.New()
	err = json.Unmarshal(payload, withdrawRequest)
	if err != nil {
		http.Error(w, "Could not decode Json", http.StatusInternalServerError)
		logger.Warn("Could not decode Json: %s", err.Error())
		return
	}
	logger = logger.GetSubLogger("withdrawId", withdrawRequest.OrderID)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = order.ValidateID(withdrawRequest.OrderID)
	if err != nil {
		logger.Warn("Invalid order format")
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}
	err = order.ValidateIDFormat(withdrawRequest.OrderID)
	if err != nil {
		logger.Warn("Invalid checksum")
		http.Error(w, "Invalid withdraw checksum", http.StatusUnprocessableEntity)
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

	if RequestContextIsClosed(ctx, w) {
		return
	}
	_, err = h.storage.GetOrder(withdrawRequest.OrderID)
	switch {
	case errors.Is(err, repositoryerrors.ErrNoContent):
	case err != nil:
		logger.Warn("error during fetching order %s", err.Error())
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "System Error", http.StatusInternalServerError)
		return
	default:
		logger.Warn("Withdraw was already registered within order database")
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "Withdraw was already registered", http.StatusConflict)
		return
	}

	if RequestContextIsClosed(ctx, w) {
		return
	}
	_, err = h.storage.GetWithdrawnWithinTransaction(ctx, withdrawRequest.OrderID)
	switch {
	case errors.Is(err, repositoryerrors.ErrNoContent):
	case err != nil:
		logger.Warn("error during fetching withdraw %s", err.Error())
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "System Error", http.StatusInternalServerError)
		return
	default:
		logger.Warn("Withdraw was already registered in withdraw database: %s", withdrawRequest.OrderID)
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "Withdraw was already registered", http.StatusConflict)
		return
	}

	logger.Debug("Trying to register withdrawn")
	logger.Debug("Fetching account info for user-account")
	if RequestContextIsClosed(ctx, w) {
		return
	}
	accountData, err := h.storage.GetCustomerAccountWithinTransaction(ctx, userID)
	if err != nil {
		logger.Warn("error during fetching account info: %s", err.Error())
		http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
		return
	}
	logger.Debug("Fetched account info=%+v", accountData)
	if accountData.Current < withdrawRequest.Sum {
		logger.Warn("error during using balance: unsufficent balance")
		http.Error(w, "unsufficent balance", http.StatusPaymentRequired)
		return
	}
	accountData.Current -= withdrawRequest.Sum
	accountData.Withdraw += withdrawRequest.Sum
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = h.storage.UpdateCustomerAccount(ctx, accountData)
	if err != nil {
		logger.Warn("error during updating account info: %s", err.Error())
		http.Error(w, "error during withDraw registration", http.StatusInternalServerError)
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = h.storage.RegisterWithdrawn(ctx, userID, withdrawn.NewRecord(withdrawRequest))
	if err != nil {
		logger.Warn("error during registering new withdrawn: %s", err.Error())
		http.Error(w, "error during withDraw registration", http.StatusInternalServerError)
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.WriteHeader(http.StatusOK)
	logger.Debug("New Withdraw is successfully registered")
}
