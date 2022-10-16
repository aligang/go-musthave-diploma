package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (h *ApiHandler) ListCustomerAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.storage.GetCustomerAccounts()
	if err != nil {
		logging.Warn("Could not get Accounts List: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	payload, err := json.Marshal(accounts)
	if err != nil {
		logging.Warn("Could not Encode JSON data: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		logging.Warn("Could set data to wire: %s", string(payload))
	}
}
