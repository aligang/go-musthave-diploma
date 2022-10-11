package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *ApiHandler) RegisterCustomerAccount(w http.ResponseWriter, r *http.Request) {
	accountInfo := &customer_account.CustomerAccount{}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Warn("Could not read data from wire")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.Unmarshal(payload, accountInfo)
	if err != nil {
		logging.Warn("Could not decode Json Data")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err = customer_account.ValidateCredentials(accountInfo); err != nil {
		logging.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	logging.Debug("Registering new account for %s", accountInfo.Login)

	var dBerr error
	h.storage.StartTransaction()
	defer func() {
		if dBerr != nil {
			h.storage.RollbackTransaction()
		}
		h.storage.CommitTransaction()
	}()

	_, err = h.storage.GetCustomerAccount(accountInfo.Login)
	if err == nil {
		logging.Warn("Account %s already exists", accountInfo.Login)
		http.Error(w, "Account %s already exists", http.StatusConflict)
	}

	err = h.storage.AddCustomerAccount(accountInfo)
	if err != nil {
		logging.Warn("Could not store Account Data: %s", err.Error())
		http.Error(w, err.Error(), http.StatusConflict)
	}
	cookie := h.auth.CreateAuthCookie(accountInfo)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	logging.Debug("account for %s is created", accountInfo.Login)
}
