package main

import (
	"context"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/tracker"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logging.Configure(os.Stdout, zerolog.DebugLevel)
	globalCtx, cancel := context.WithCancel(context.Background())
	cfg := config.Init()
	Storage := storage.New(cfg)
	Auth := auth.New()
	Tracker := tracker.New(Storage, cfg)
	Tracker.RunInBackground(globalCtx)
	mux := handler.New(Storage, Auth, cfg)
	mux.ApplyProdConfig()
	//mux.ApplyDebugConfig()

	app := New(globalCtx, mux, cfg)

	go runServer(app)
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	cancel()
	logging.Debug("Server stopped")
}

func runServer(server *http.Server) {
	logging.Debug("enable TCP listener on: %s", server.Addr)
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}
	logging.Debug(" Starting Server on: %s", server.Addr)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

type contextKeyType struct {
	Key string
}

type contextValueType struct {
	Value string
}

var contextKey = contextKeyType{"server-context"}
var contextValue = contextValueType{""}

func New(ctx context.Context, mux http.Handler, cfg *config.Config) *http.Server {
	serverBaseCtxFunc := func(listener net.Listener) context.Context {
		return context.WithValue(ctx, "context-name", contextValue)
	}
	return &http.Server{
		Addr:        cfg.RunAddress,
		Handler:     mux,
		BaseContext: serverBaseCtxFunc,
	}
}
