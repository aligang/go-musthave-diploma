package storage

import (
	"context"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/database"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
)
import "github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"

type Test interface {
	Check(int, int)
}

type Storage interface {
	StartTransaction(ctx context.Context)
	RollbackTransaction(ctx context.Context)
	CommitTransaction(ctx context.Context)

	AddCustomerAccount(ctx context.Context, customerAccount *customer_account.CustomerAccount) error
	UpdateCustomerAccount(ctx context.Context, customerAccount *customer_account.CustomerAccount) error
	GetCustomerAccount(login string) (*customer_account.CustomerAccount, error)
	GetCustomerAccountWithinTransaction(ctx context.Context, login string) (*customer_account.CustomerAccount, error)
	GetCustomerAccounts() (customer_account.CustomerAccounts, error)

	AddOrder(ctx context.Context, userId string, order *order.Order) error
	GetOrder(orderId string) (*order.Order, error)
	GetOrderWithinTransaction(ctx context.Context, orderId string) (*order.Order, error)
	ListOrders(userId string) ([]order.Order, error)
	GetOrderOwner(ctx context.Context, orderId string) (string, error)
	UpdateOrder(ctx context.Context, order *order.Order) error

	AddOrderToPendingList(ctx context.Context, orderId string) error
	GetPendingOrders(ctx context.Context) ([]string, error)
	RemoveOrderFromPendingList(ctx context.Context, orderId string) error

	RegisterWithdrawn(ctx context.Context, userId string, withdraw *withdrawn.WithdrawnRecord) error
	GetWithdrawnWithinTransaction(ctx context.Context, orderId string) (*withdrawn.WithdrawnRecord, error)
	ListWithdrawns(orderId string) ([]withdrawn.WithdrawnRecord, error)
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
