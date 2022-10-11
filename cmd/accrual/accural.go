package main

import (
	"github.com/aligang/go-musthave-diploma/internal/accural/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)
import "github.com/aligang/go-musthave-diploma/internal/accural/handler"

var URI = "127.0.0.1:9090"

func main() {
	logging.Configure(os.Stdout, zerolog.DebugLevel)
	storage := memory.New()
	mux := handler.New(storage)

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)

	mux.Get("/api/orders/{order}", mux.Fetch)
	mux.Get("/api/orders", mux.BulkFetch)
	mux.Post("/api/orders/", mux.Update)
	logging.Debug(" Starting Server on: %s", URI)
	logging.Crit(http.ListenAndServe(URI, mux).Error())

}
