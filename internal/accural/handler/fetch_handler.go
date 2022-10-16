package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *APIHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order")
	order, err := h.Storage.Get(orderID)
	if err != nil {
		logging.Warn("Could not found record with id=%s", orderID)
		http.Error(w, "Record was not found", http.StatusNotFound)
		return
	}
	j, err := json.Marshal(&order)
	if err != nil {
		logging.Warn("Could not encode JSON data")
		http.Error(w, "Could not encode JSON data", http.StatusUnsupportedMediaType)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	logging.Debug("%s Sending Storage", string(j))
	_, err = w.Write(j)
	if err != nil {
		logging.Warn("Error during data transfer", orderID)
		http.Error(w, "Record was not found", http.StatusUnsupportedMediaType)
		return
	}
}
