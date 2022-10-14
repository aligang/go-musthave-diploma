package handler

import (
	"encoding/json"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repository_errors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *ApiHandler) RegisterCustomerAccount(w http.ResponseWriter, r *http.Request) {
	//r.Context()
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
	logging.Debug("Registering new account for %s", accountInfo.Login)

	var dBerr error
	h.storage.StartTransaction()
	defer func() {
		if dBerr != nil {
			h.storage.RollbackTransaction()
		}
		h.storage.CommitTransaction()
	}()

	_, err = h.storage.GetCustomerAccountWithinTransaction(accountInfo.Login)
	switch {
	case errors.Is(err, repository_errors.ErrNoContent):
	case err != nil:
		logging.Warn("error during fetching Account %s", err.Error())
		http.Error(w, "Account %s already exists", http.StatusInternalServerError)
		return
	default:
		logging.Warn("Account %s already exists", accountInfo.Login)
		http.Error(w, "Account %s already exists", http.StatusConflict)
		return
	}

	err = h.storage.AddCustomerAccount(accountInfo)
	if err != nil {
		logging.Warn("Could not store Account Data: %s", err.Error())
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	cookie := h.auth.CreateAuthCookie(accountInfo)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	logging.Debug("account for %s is created", accountInfo.Login)
}

//
