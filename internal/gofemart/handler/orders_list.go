package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
	"sort"
)

func (h *APIhandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	logging.Warn("Processing order list request")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	userID, err := auth.ResolveUsername(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logging.Warn("No user info were provided")
		return
	}

	h.storage.StartTransaction(ctx)
	defer func() {
		h.storage.CommitTransaction(ctx)
	}()

	logging.Debug("Fetching orders registered for user=%s from repository", userID)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	orders, err := h.storage.ListOrders(userID)
	if err != nil {
		http.Error(w, "error during Fetching orders", http.StatusInternalServerError)
		logging.Warn("Error during fetching orders register for user=%s: %s", userID, err.Error())
		return
	} else if len(orders) == 0 {
		http.Error(w, "there is no registered orders", http.StatusNoContent)
		logging.Warn("User=%s has no registered orders", userID)
		return
	}

	logging.Debug("user %s, has registered orders: %+v", userID, orders)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	sort.Sort(order.OrderSlice(orders))
	ordersPayload, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "error during Fetching orders", http.StatusInternalServerError)
		logging.Warn("Could not decode json")
		return
	}
	logging.Debug("forming response %s", string(ordersPayload))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersPayload)
	if err != nil {
		logging.Debug("Error during writing data to wire")
	}
	logging.Debug("orders  list user=%s was sent", userID)
}
