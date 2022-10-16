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
	logging.Warn("Processing withdraw list request")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	userID, err := auth.ResolveUsername(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logging.Warn("No user info were provided")
		return
	}
	logging.Debug("Processing Withdraw for user %s", userID)
	payload, err := io.ReadAll(r.Body)
	logging.Warn("Withdraw add: received request %s", string(payload))
	if err != nil {
		http.Error(w, "Could not read data from wire", http.StatusInternalServerError)
		return
	}

	withdrawRequest := withdrawn.New()
	err = json.Unmarshal(payload, withdrawRequest)
	if err != nil {
		http.Error(w, "Could not decode Json", http.StatusInternalServerError)
		logging.Warn("Could not decode Json: %s", err.Error())
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = order.ValidateID(withdrawRequest.Order)
	if err != nil {
		logging.Warn("Invalid order format: %s", withdrawRequest.Order)
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}
	err = order.ValidateIDFormat(withdrawRequest.Order)
	if err != nil {
		logging.Warn("Invalid orderID checksum: %s", withdrawRequest.Order)
		http.Error(w, "Invalid orderID checksum", http.StatusUnprocessableEntity)
		return
	}
	logging.Debug("Withdraw ID is %s", withdrawRequest.Order)

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
	_, err = h.storage.GetOrder(withdrawRequest.Order)
	switch {
	case errors.Is(err, repositoryerrors.ErrNoContent):
	case err != nil:
		logging.Warn("error during fetching order %s", err.Error())
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "System Error", http.StatusInternalServerError)
		return
	default:
		logging.Warn("Withdraw=%s was already registered within order database", withdrawRequest.Order)
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "Withdraw was already registered", http.StatusConflict)
		return
	}

	if RequestContextIsClosed(ctx, w) {
		return
	}
	_, err = h.storage.GetWithdrawnWithinTransaction(ctx, withdrawRequest.Order)
	switch {
	case errors.Is(err, repositoryerrors.ErrNoContent):
	case err != nil:
		logging.Warn("error during fetching withdraw %s", err.Error())
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "System Error", http.StatusInternalServerError)
		return
	default:
		logging.Warn("Withdraw was already registered in withdraw database: %s", withdrawRequest.Order)
		if RequestContextIsClosed(ctx, w) {
			return
		}
		http.Error(w, "Withdraw was already registered", http.StatusConflict)
		return
	}

	logging.Debug("Trying to register withdrawn order=%s to user-account=%s", withdrawRequest.Order, userID)
	logging.Debug("Fetching account info for user-account=%s", userID)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	accountData, err := h.storage.GetCustomerAccountWithinTransaction(ctx, userID)
	if err != nil {
		logging.Warn("error during fetching account info: %s", err.Error())
		http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
		return
	}
	logging.Debug("Fetched account info=%+v", accountData)
	if accountData.Current < withdrawRequest.Sum {
		logging.Warn("error during using balance of: %s, unsufficent balance", userID)
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
		logging.Warn("error during updating account info: %s", err.Error())
		http.Error(w, "error during withDraw registration", http.StatusInternalServerError)
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = h.storage.RegisterWithdrawn(ctx, userID, withdrawn.NewRecord(withdrawRequest))
	if err != nil {
		logging.Warn("error during registering new withdrawn=%s  for account=%s: %s",
			withdrawRequest.Order, accountData.Login, err.Error())
		http.Error(w, "error during withDraw registration", http.StatusInternalServerError)
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.WriteHeader(http.StatusOK)
	logging.Debug("New Withdraw=%s is successfully registered", withdrawRequest.Order)
}
