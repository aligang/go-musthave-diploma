package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	if err = customer_account.ValidateCredentials(accountInfo); err != nil {
		logging.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logging.Debug("Authenticating user %s", accountInfo.Login)
	account, err := h.storage.GetCustomerAccount(accountInfo.Login)
	if err != nil || accountInfo.Password != account.Password {
		http.Error(w, "Authentication Failure", http.StatusUnauthorized)
		logging.Debug("Could not authenticate user %s", accountInfo.Login)
		return
	}

	cookie := h.auth.CreateAuthCookie(accountInfo)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	logging.Debug("%s is authenticated", accountInfo.Login)
}
