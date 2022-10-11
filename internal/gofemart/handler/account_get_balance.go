package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (h *ApiHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	logging.Warn("Processing account balance request")
	userId, err := auth.ResolveUsername(r)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		logging.Warn("No user info were provided")
		return
	}
	accountInfo, err := h.storage.GetCustomerAccount(userId)
	if err != nil {
		http.Error(w, "Could not provide account balance info", http.StatusInternalServerError)
		logging.Warn("Error during fetching account info from repository")
		return
	}
	payload, err := json.Marshal(accountInfo.AccountBalance)
	if err != nil {
		http.Error(w, "Could not provide account balance info", http.StatusInternalServerError)
		logging.Warn("Could not encode %+v", accountInfo.AccountBalance)
		return
	}
	logging.Debug("Sending data to wire", string(payload))
	w.Header().Set("Content-Type", "Application/Json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)
	if err != nil {
		logging.Warn("Could not write data to wire")
		return
	}
	logging.Warn("Account balance response successfully sent")
}
