package handler

import (
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage"
	"github.com/go-chi/chi/v5"
)

type ApiHandler struct {
	*chi.Mux
	storage storage.Storage
	auth    *auth.Auth
	config  *config.Config
}

func New(storage storage.Storage, auth *auth.Auth, cfg *config.Config) *ApiHandler {
	return &ApiHandler{
		Mux:     chi.NewMux(),
		storage: storage,
		auth:    auth,
		config:  cfg,
	}
}
