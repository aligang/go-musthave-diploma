package handler

import (
	"encoding/json"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *APIhandler) RegisterCustomerAccount(w http.ResponseWriter, r *http.Request) {
	logger := logging.Logger.GetSubLogger("Method", "RegisterAccount")
	logger.Debug("Processing request")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	accountInfo := account.New()
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
	if err = account.ValidateCredentials(accountInfo); err != nil {
		logging.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.GetSubLogger("userID", "accountInfo.Login")
	logger.Debug("Registering new account")
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
	_, err = h.storage.GetCustomerAccountWithinTransaction(ctx, accountInfo.Login)
	switch {
	case errors.Is(err, repositoryerrors.ErrNoContent):
	case err != nil:
		logger.Warn("error during fetching Account %s", err.Error())
		http.Error(w, "Account %s already exists", http.StatusInternalServerError)
		return
	default:
		logger.Warn("Account %s already exists", accountInfo.Login)
		http.Error(w, "Account %s already exists", http.StatusConflict)
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	err = h.storage.AddCustomerAccount(ctx, accountInfo)
	if err != nil {
		logger.Warn("Could not store Account Data: %s", err.Error())
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if RequestContextIsClosed(ctx, w) {
		return
	}
	cookie := h.auth.CreateAuthCookie(accountInfo)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	logger.Debug("account for %s is created", accountInfo.Login)
}

//
