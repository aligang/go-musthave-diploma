package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (h *ApiHandler) ListPendingOrders(w http.ResponseWriter, r *http.Request) {
	logging.Debug("pending orders  list was requested")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	orderIds, err := h.storage.GetPendingOrders(ctx)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logging.Warn("Error during fetching pending orders list")
		return
	}
	logging.Debug("user %s, has registered orders")
	if RequestContextIsClosed(ctx, w) {
		return
	}
	ordersPayload, err := json.Marshal(orderIds)
	logging.Debug("forming response %s", string(ordersPayload))
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersPayload)
	if err != nil {
		logging.Debug("Error during writing data to wire")
	}
	logging.Debug("pending orders  list was sent")
}
