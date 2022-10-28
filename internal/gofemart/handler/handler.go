package handler

import (
	"context"
	"github.com/aligang/go-musthave-diploma/internal/accural"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/compress"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type APIhandler struct {
	*chi.Mux
	storage       storage.Storage
	auth          *auth.Auth
	config        *config.Config
	accrualClient *accural.AccrualClient
}

func New(storage storage.Storage, auth *auth.Auth, cfg *config.Config) *APIhandler {
	return &APIhandler{
		Mux:           chi.NewMux(),
		storage:       storage,
		auth:          auth,
		config:        cfg,
		accrualClient: accural.New(cfg),
	}
}

func (h *APIhandler) ApplyProdConfig() {
	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Recoverer)
	//user api
	h.With(compress.Gzip).
		Post("/api/user/register", h.RegisterCustomerAccount)
	h.With(compress.Gzip).
		Post("/api/user/login", h.Login)
	h.With(h.auth.CheckAuthInfo).
		With(compress.Gzip).
		Post("/api/user/orders", h.AddOrder)
	h.With(h.auth.CheckAuthInfo).
		With(compress.Gzip).
		Get("/api/user/orders", h.ListOrders)
	h.With(h.auth.CheckAuthInfo).
		With(compress.Gzip).
		Get("/api/user/balance", h.GetAccountBalance)
	h.With(h.auth.CheckAuthInfo).
		With(compress.Gzip).
		Post("/api/user/balance/withdraw", h.AddWithdraw)
	h.With(h.auth.CheckAuthInfo).
		With(compress.Gzip).
		Get("/api/user/withdrawals", h.ListWithdraws)
}

func (h *APIhandler) ApplyDebugConfig() {
	h.Get("/api/internal/accounts", h.ListCustomerAccounts)
}

func RequestContextIsClosed(ctx context.Context, w http.ResponseWriter) bool {
	select {
	default:
		return false
	case <-ctx.Done():
		logging.Debug("Context was cancelled")
		http.Error(w, "Context was cancelled", http.StatusInternalServerError)
		return true
	}
}
