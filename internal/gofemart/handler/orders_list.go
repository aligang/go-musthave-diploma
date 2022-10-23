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
	logger := logging.Logger.GetSubLogger("Method", "Order List")
	logger.Warn("Processing request")
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

	logger = logger.GetSubLogger("userID", userID)
	h.storage.StartTransaction(ctx)
	defer func() {
		h.storage.CommitTransaction(ctx)
	}()

	logger.Debug("Fetching orders registered  from repository")
	if RequestContextIsClosed(ctx, w) {
		return
	}
	orders, err := h.storage.ListOrders(userID)
	if err != nil {
		http.Error(w, "error during Fetching orders", http.StatusInternalServerError)
		logger.Warn("Error during fetching orders register for user: %s", err.Error())
		return
	} else if len(orders) == 0 {
		http.Error(w, "there is no registered orders", http.StatusNoContent)
		logger.Warn("User has no registered orders")
		return
	}

	logger.Debug("user %s, has registered orders: %+v", userID, orders)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	sort.Sort(order.OrderSlice(orders))
	ordersPayload, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "error during Fetching orders", http.StatusInternalServerError)
		logger.Warn("Could not decode json")
		return
	}
	logger.Debug("forming response %s", string(ordersPayload))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersPayload)
	if err != nil {
		logger.Debug("Error during writing data to wire")
	}
	logger.Debug("orders  list was sent")
}
