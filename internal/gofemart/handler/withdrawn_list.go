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
	logger := logging.Logger.GetSubLogger("Method", "Withdrawn List")
	logger.Warn("Processing request")
	ctx := r.Context()
	if RequestContextIsClosed(ctx, w) {
		return
	}
	userID, err := auth.ResolveUsername(r)
	logger = logger.GetSubLogger("userID", userID)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		logger.Warn("No user info were provided")
		return
	}

	if RequestContextIsClosed(ctx, w) {
		return
	}
	h.storage.StartTransaction(ctx)
	defer func() {
		h.storage.CommitTransaction(ctx)
	}()

	logger.Debug("Fetching withdraws  from repository")
	if RequestContextIsClosed(ctx, w) {
		return
	}
	withdrawns, err := h.storage.ListWithdrawns(userID)
	switch {
	case err != nil:
		http.Error(w, "", http.StatusInternalServerError)
		logger.Warn("Error during fetching withdraws registered for user: %s", err.Error())
		return
	case len(withdrawns) == 0:
		http.Error(w, "there is now registered withdraws", http.StatusNoContent)
		logger.Warn("User=%s has now registered withdraws", userID)
		return
	}

	logger.Debug("user has registered withdrawns: %+v", withdrawns)
	sort.Sort(withdrawn.WithdrawnSlice(withdrawns))
	withdrawsPayload, err := json.Marshal(withdrawns)
	if err != nil {
		http.Error(w, "error during Fetching orders", http.StatusInternalServerError)
		logger.Warn("Could not decode json")
		return
	}
	logger.Debug("forming response %s", string(withdrawsPayload))
	if RequestContextIsClosed(ctx, w) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(withdrawsPayload)
	if err != nil {
		logger.Debug("Error during writing data to wire")
	}
	logger.Debug("orders list was sent")
}
