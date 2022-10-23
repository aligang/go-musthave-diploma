package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (h *APIhandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	logger := logging.Logger.GetSubLogger("Method", "GetAccount Balance")
	logging.Warn("Processing request")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	userID, err := auth.ResolveUsername(r)
	logger = logger.GetSubLogger("userID", userID)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		logger.Warn("No user info were provided")
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	accountInfo, err := h.storage.GetCustomerAccount(userID)
	if err != nil {
		http.Error(w, "Could not provide account balance info", http.StatusInternalServerError)
		logger.Warn("Error during fetching account info from repository")
		return
	}
	payload, err := json.Marshal(accountInfo.AccountBalance)
	if err != nil {
		http.Error(w, "Could not provide account balance info", http.StatusInternalServerError)
		logger.Warn("Could not encode %+v", accountInfo.AccountBalance)
		return
	}
	logger.Debug("Sending data to wire %s", string(payload))
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)
	if err != nil {
		logger.Warn("Could not write data to wire")
		return
	}
	logger.Warn("Account balance response successfully sent")
}
