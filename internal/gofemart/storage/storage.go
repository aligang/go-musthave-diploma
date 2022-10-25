package storage

import (
	"context"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/database"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"github.com/jmoiron/sqlx"
)
import "github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"

type Test interface {
	Check(int, int)
}

type Storage interface {
	WithinTransaction(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error

	AddCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount, tx *sqlx.Tx) error
	UpdateCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount, tx *sqlx.Tx) error
	GetCustomerAccount(ctx context.Context, login string, tx *sqlx.Tx) (*account.CustomerAccount, error)
	GetCustomerAccounts() (account.CustomerAccounts, error)

	AddOrder(ctx context.Context, userID string, order *order.Order, tx *sqlx.Tx) error
	GetOrder(ctx context.Context, orderID string, tx *sqlx.Tx) (*order.Order, error)
	ListOrders(ctx context.Context, userID string) ([]order.Order, error)
	GetOrderOwner(ctx context.Context, orderID string, tx *sqlx.Tx) (string, error)
	UpdateOrder(ctx context.Context, order *order.Order, tx *sqlx.Tx) error

	AddOrderToPendingList(ctx context.Context, orderID string, tx *sqlx.Tx) error
	GetPendingOrders(ctx context.Context, tx *sqlx.Tx) ([]string, error)
	RemoveOrderFromPendingList(ctx context.Context, orderID string, tx *sqlx.Tx) error

	RegisterWithdrawn(ctx context.Context, userID string, withdraw *withdrawn.WithdrawnRecord, tx *sqlx.Tx) error
	GetWithdrawn(ctx context.Context, orderID string, tx *sqlx.Tx) (*withdrawn.WithdrawnRecord, error)
	ListWithdrawns(ctx context.Context, orderID string) ([]withdrawn.WithdrawnRecord, error)
}

func New(config *config.Config) Storage {
	var storage Storage
	logging.Debug("Initialization Storage")
	if len(config.DatabaseURI) == 0 {
		storage = memory.New()
	} else {
		storage = database.New(config)
	}
	logging.Debug("Storage Initialization finished")
	return storage
}
