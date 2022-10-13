package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"net/http"
	"sort"
)

func (h *ApiHandler) ListWithdraws(w http.ResponseWriter, r *http.Request) {
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

	logging.Debug("Fetching withdraws registered for user=%s from repository", userId)
	withdrawns, err := h.storage.ListWithdrawns(userId)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logging.Warn("Error during fetching withdraws registered for user=%s: %s", userId, err.Error())
		return
	}
	if len(withdrawns) == 0 {
		http.Error(w, "there is now registered withdraws", http.StatusNoContent)
		logging.Warn("User=%s has now registered withdraws", userId)
		return
	}

	logging.Debug("user %s, has registered orders: %+v", userId, withdrawns)
	sort.Sort(withdrawn.WithdrawnSlice(withdrawns))
	withdrawsPayload, err := json.Marshal(withdrawns)
	logging.Debug("forming response %s", string(withdrawsPayload))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(withdrawsPayload)
	if err != nil {
		logging.Debug("Error during writing data to wire")
	}
	logging.Debug("orders  list user=%s was sent", userId)
}
