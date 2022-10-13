package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (h *ApiHandler) ListPendingOrders(w http.ResponseWriter, r *http.Request) {

	logging.Debug("pending orders  list was requested")
	orderIds, err := h.storage.GetPendingOrders()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logging.Warn("Error during fetching pending orders list")
		return
	}
	logging.Debug("user %s, has registered orders")

	ordersPayload, err := json.Marshal(orderIds)
	logging.Debug("forming response %s", string(ordersPayload))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersPayload)
	if err != nil {
		logging.Debug("Error during writing data to wire")
	}
	logging.Debug("pending orders  list was sent")
}
