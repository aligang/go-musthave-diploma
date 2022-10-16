package storage

import (
	"context"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
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

	AddCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount) error
	UpdateCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount) error
	GetCustomerAccount(login string) (*account.CustomerAccount, error)
	GetCustomerAccountWithinTransaction(ctx context.Context, login string) (*account.CustomerAccount, error)
	GetCustomerAccounts() (account.CustomerAccounts, error)

	AddOrder(ctx context.Context, userId string, order *order.Order) error
	GetOrder(orderID string) (*order.Order, error)
	GetOrderWithinTransaction(ctx context.Context, orderID string) (*order.Order, error)
	ListOrders(userId string) ([]order.Order, error)
	GetOrderOwner(ctx context.Context, orderID string) (string, error)
	UpdateOrder(ctx context.Context, order *order.Order) error

	AddOrderToPendingList(ctx context.Context, orderID string) error
	GetPendingOrders(ctx context.Context) ([]string, error)
	RemoveOrderFromPendingList(ctx context.Context, orderID string) error

	RegisterWithdrawn(ctx context.Context, userId string, withdraw *withdrawn.WithdrawnRecord) error
	GetWithdrawnWithinTransaction(ctx context.Context, orderID string) (*withdrawn.WithdrawnRecord, error)
	ListWithdrawns(orderID string) ([]withdrawn.WithdrawnRecord, error)
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
