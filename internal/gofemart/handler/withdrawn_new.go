package handler

import (
	"encoding/json"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repository_errors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"io"
	"net/http"
)

func (h *ApiHandler) AddWithdraw(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.ResolveUsername(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logging.Warn("No user info were provided")
		return
	}
	logging.Debug("Processing Withdraw for user %s", userId)
	payload, err := io.ReadAll(r.Body)
	logging.Warn("Withdraw add: received request %s", string(payload))
	if err != nil {
		http.Error(w, "Could not read data from wire", http.StatusInternalServerError)
		return
	}

	withdrawRequest := &withdrawn.Withdrawn{}
	err = json.Unmarshal(payload, withdrawRequest)
	if err != nil {
		http.Error(w, "Could not decode Json", http.StatusInternalServerError)
		logging.Warn("Could not decode Json: %s", err.Error())
		return
	}
	err = order.ValidateId(withdrawRequest.Order)
	if err != nil {
		logging.Warn("Invalid order format: %s", withdrawRequest.Order)
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}
	err = order.ValidateIdFormat(withdrawRequest.Order)
	if err != nil {
		logging.Warn("Invalid orderId checksum: %s", withdrawRequest.Order)
		http.Error(w, "Invalid orderId checksum", http.StatusUnprocessableEntity)
		return
	}
	logging.Debug("Withdraw ID is %d", withdrawRequest.Order)
	var dBerr error
	h.storage.StartTransaction()
	defer func() {
		if dBerr != nil {
			h.storage.RollbackTransaction()
		}
		h.storage.CommitTransaction()
	}()

	_, err = h.storage.GetOrder(withdrawRequest.Order)
	switch {
	case errors.Is(err, repository_errors.ErrNoContent):
	case err != nil:
		logging.Warn("error during fetching order %s", err.Error())
		http.Error(w, "System Error", http.StatusInternalServerError)
		return
	default:
		logging.Warn("Withdraw=%s was already registered within order database", withdrawRequest.Order)
		http.Error(w, "Withdraw was already registered", http.StatusConflict)
		return
	}

	_, err = h.storage.GetWithdrawnWithinTransaction(withdrawRequest.Order)
	switch {
	case errors.Is(err, repository_errors.ErrNoContent):
	case err != nil:
		logging.Warn("error during fetching withdraw %s", err.Error())
		http.Error(w, "System Error", http.StatusInternalServerError)
		return
	default:
		logging.Warn("Withdraw was already registered in withdraw database", withdrawRequest.Order)
		http.Error(w, "Withdraw was already registered", http.StatusConflict)
		return
	}

	logging.Debug("Trying to register withdrawn order=%s to user-account=%s", withdrawRequest.Order, userId)
	logging.Debug("Fetching account info for user-account=%s", userId)
	accountData, err := h.storage.GetCustomerAccountWithinTransaction(userId)
	if err != nil {
		logging.Warn("error during fetching account info: %s", err.Error())
		http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
		return
	}
	logging.Debug("Fetched account info=%+v", accountData)
	if accountData.Current < withdrawRequest.Sum {
		logging.Warn("error during using balance of: %s, unsufficent balance", userId)
		http.Error(w, "unsufficent balance", http.StatusPaymentRequired)
		return
	}
	accountData.Current -= withdrawRequest.Sum
	accountData.Withdraw += withdrawRequest.Sum
	err = h.storage.UpdateCustomerAccount(accountData)
	if err != nil {
		logging.Warn("error during updating account info: %s", err.Error())
		http.Error(w, "error during withDraw registration", http.StatusInternalServerError)
		return
	}
	err = h.storage.RegisterWithdrawn(userId, withdrawn.NewRecord(withdrawRequest))
	if err != nil {
		logging.Warn("error during registering new withdrawn=%s  for account=%s: %s",
			withdrawRequest.Order, accountData.Login, err.Error())
		http.Error(w, "error during withDraw registration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logging.Debug("New Withdraw=%s is successfully registered", withdrawRequest.Order)
}
