package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
)

func (h *APIHandler) Update(w http.ResponseWriter, r *http.Request) {
	m := message.AccuralMessage{}
	payload, err := io.ReadAll(r.Body)
	logging.Debug("Recieved JSON: %s", string(payload))
	err = json.Unmarshal(payload, &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()
	if err != nil {
		logging.Crit("Could not decode received JSON")
		http.Error(w, "Mailformed JSON", http.StatusBadRequest)
		return
	}

	err = h.Storage.Put(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
