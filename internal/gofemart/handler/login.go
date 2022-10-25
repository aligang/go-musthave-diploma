package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *APIhandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := logging.Logger.GetSubLogger("Method", "login")
	logger.Warn("Processing request")
	ctx := r.Context()
	accountInfo := account.New()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Warn("Could not read data from wire")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.Unmarshal(payload, accountInfo)
	if err != nil {
		logger.Warn("Could not decode Json Data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = account.ValidateCredentials(accountInfo); err != nil {
		logger.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.GetSubLogger("userID", accountInfo.Login)
	logger.Debug("Authenticating user")
	account, err := h.storage.GetCustomerAccount(ctx, accountInfo.Login, nil)
	if err != nil || accountInfo.Password != account.Password {
		http.Error(w, "Authentication Failure", http.StatusUnauthorized)
		logger.Debug("Could not authenticate user")
		return
	}
	cookie := h.auth.CreateAuthCookie(accountInfo)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	logger.Debug("login request successfully processed")
}
