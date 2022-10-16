package handler

import (
	"encoding/json"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"net/http"
	"sort"
)

func (h *APIhandler) ListWithdraws(w http.ResponseWriter, r *http.Request) {
	logging.Warn("Processing withdraw list request")
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

	if RequestContextIsClosed(ctx, w) {
		return
	}
	h.storage.StartTransaction(ctx)
	defer func() {
		h.storage.CommitTransaction(ctx)
	}()

	logging.Debug("Fetching withdraws registered for user=%s from repository", userID)
	if RequestContextIsClosed(ctx, w) {
		return
	}
	withdrawns, err := h.storage.ListWithdrawns(userID)
	switch {
	case err != nil:
		http.Error(w, "", http.StatusInternalServerError)
		logging.Warn("Error during fetching withdraws registered for user=%s: %s", userID, err.Error())
		return
	case len(withdrawns) == 0:
		http.Error(w, "there is now registered withdraws", http.StatusNoContent)
		logging.Warn("User=%s has now registered withdraws", userID)
		return
	}

	logging.Debug("user %s, has registered orders: %+v", userID, withdrawns)
	sort.Sort(withdrawn.WithdrawnSlice(withdrawns))
	withdrawsPayload, err := json.Marshal(withdrawns)
	if err != nil {
		http.Error(w, "error during Fetching orders", http.StatusInternalServerError)
		logging.Warn("Could not decode json")
		return
	}
	logging.Debug("forming response %s", string(withdrawsPayload))
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(withdrawsPayload)
	if err != nil {
		logging.Debug("Error during writing data to wire")
	}
	logging.Debug("orders  list user=%s was sent", userID)
}
