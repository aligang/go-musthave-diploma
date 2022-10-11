package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (h *ApiHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.ResolveUsername(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logging.Warn("No user info were provided")
		return
	}

	h.storage.StartTransaction()
	defer func() {
		h.storage.CommitTransaction()
	}()

	logging.Debug("Fetching orders registered for user=%s from repository", userId)
	orders, err := h.storage.ListOrders(userId)

	if err != nil {
		http.Error(w, "error during Fetching orders", http.StatusInternalServerError)
		logging.Warn("Error during fetching orders register for user=%s: %s", userId, err.Error())
		return
	}
	if len(orders) == 0 {
		http.Error(w, "there is now registered orders", http.StatusNoContent)
		logging.Warn("User=%s has now registered orders", userId)
		return
	}

	logging.Debug("user %s, has registered orders: %+v", userId, orders)
	order.Sort(orders)
	ordersPayload, err := json.Marshal(orders)
	logging.Debug("forming response %s", string(ordersPayload))
	w.Header().Set("Content-Type", "Application/Json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersPayload)
	if err != nil {
		logging.Debug("Error during writing data to wire")
	}
	logging.Debug("orders  list user=%s was sent", userId)
}
