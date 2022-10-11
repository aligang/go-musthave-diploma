package main

import (
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/tracker"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

var URI = "127.0.0.1:8080"

func main() {
	logging.Configure(os.Stdout, zerolog.DebugLevel)
	cfg := config.Init()
	Storage := storage.New(cfg)
	Auth := auth.New()
	tracker.New(Storage, cfg)
	mux := handler.New(Storage, Auth, cfg)

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	//user api
	mux.Post("/api/user/register", mux.RegisterCustomerAccount)
	mux.Post("/api/user/login", mux.Login)
	mux.With(Auth.CheckAuthInfo).
		Post("/api/user/orders", mux.AddOrder)
	mux.With(Auth.CheckAuthInfo).
		Get("/api/user/orders", mux.ListOrders)
	mux.With(Auth.CheckAuthInfo).
		Get("/api/user/balance", mux.GetAccountBalance)
	mux.With(Auth.CheckAuthInfo).
		Post("/api/user/balance/withdrawn", mux.AddWithdraw)
	mux.With(Auth.CheckAuthInfo).
		Get("/api/user/withdrawals", mux.ListWithdraws)

	//internal api
	mux.Get("/api/internal/accounts", mux.ListCustomerAccounts)
	mux.Get("/api/internal/pending-orders", mux.ListPendingOrders)

	logging.Debug(" Starting Server on: %s", URI)
	logging.Crit(http.ListenAndServe(URI, mux).Error())
}
