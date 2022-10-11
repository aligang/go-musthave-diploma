package handler

import "github.com/go-chi/chi/v5"
import "github.com/aligang/go-musthave-diploma/internal/accural/storage/memory"

type APIHandler struct {
	*chi.Mux
	Storage *memory.Storage
	//Config  Config
}

func New(storage *memory.Storage) APIHandler {
	mux := APIHandler{
		Mux:     chi.NewMux(),
		Storage: storage,
	}
	return mux
}
