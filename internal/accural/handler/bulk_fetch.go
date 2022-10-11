package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (h *APIHandler) BulkFetch(w http.ResponseWriter, r *http.Request) {
	orders := h.Storage.BulkGet()

	j, err := json.Marshal(&orders)
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
		logging.Warn("Error during data transfer")
		http.Error(w, "Error during data transfer", http.StatusInternalServerError)
		return
	}
}
