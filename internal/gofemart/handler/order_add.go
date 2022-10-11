package handler

import (
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/accural"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order/status"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *ApiHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.ResolveUsername(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logging.Warn("No user info were provided")
		return
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read data from wire", http.StatusInternalServerError)
		return
	}
	orderId := string(payload)
	err = order.ValidateId(orderId)
	if err != nil {
		logging.Warn("Invalid order format: %s", orderId)
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}
	err = order.ValidateIdFormat(orderId)
	if err != nil {
		logging.Warn("Invalid orderId checksum: %s", orderId)
		http.Error(w, "Invalid orderId checksum", http.StatusUnprocessableEntity)
		return
	}

	var dBerr error
	h.storage.StartTransaction()
	defer func() {
		if dBerr != nil {
			h.storage.RollbackTransaction()
		}
		h.storage.CommitTransaction()
	}()

	_, err = h.storage.GetOrder(orderId)
	if err == nil {
		owner, err := h.storage.GetOrderOwner(orderId)
		if owner == userId {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			return
		}
		if err != nil {
			logging.Warn("Could not check owner for Order=%s", orderId)
			http.Error(w, "error during registering order", http.StatusInternalServerError)
			return
		}
		logging.Warn("Order=%s was already registered in order database", orderId)
		http.Error(w, "Order was already registered", http.StatusConflict)
		return
	}
	_, err = h.storage.GetWithdrawn(orderId)
	if err == nil {
		logging.Warn("Order=%s was already registered in order database", orderId)
		http.Error(w, "Order was already registered", http.StatusInternalServerError)
		return
	}

	logging.Warn("Trying to bind order=%s to user-account=%s", orderId, userId)
	accuralRecord, err := accural.FetchOrderInfo(orderId, h.config)
	if err != nil {
		logging.Warn("Failed to fecth accural Info: %s", err.Error())
		http.Error(w, "error during registering order", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	err = h.storage.AddOrder(userId, accuralRecord)
	if err != nil {
		logging.Warn("error during adding order to order Database: %s", err.Error())
		http.Error(w, "error during registering order", http.StatusInternalServerError)
		return
	}
	if status.RequiresTracking(accuralRecord.Status) {
		logging.Debug("adding order=%s to PENDING", accuralRecord.Order)
		err = h.storage.AddOrderToPendingList(orderId)
		if err != nil {
			logging.Warn("error during adding order to pending list: %s", err.Error())
			http.Error(w, "error during registering order", http.StatusInternalServerError)
			return
		}
	}
	if accuralRecord.Status == status.PROCESSED {
		logging.Debug("Applying orderId=%s accural to %s balance", accuralRecord.Order, userId)
		accountData, err := h.storage.GetCustomerAccount(userId)
		if err != nil {
			logging.Warn("error during fetching account info: %s", err.Error())
			http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
			return
		}
		accountData.Balance += accuralRecord.Accural
		err = h.storage.UpdateCustomerAccount(accountData)
		if err != nil {
			logging.Warn("error during fetching account info: %s", err.Error())
			http.Error(w, "error during add accural to balance", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	logging.Debug("New Order=%s is successfully registered", orderId)
}
