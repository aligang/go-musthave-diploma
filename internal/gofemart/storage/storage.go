package storage

import (
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
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
	StartTransaction()
	RollbackTransaction()
	CommitTransaction()

	AddCustomerAccount(customerAccount *customer_account.CustomerAccount) error
	UpdateCustomerAccount(customerAccount *customer_account.CustomerAccount) error
	GetCustomerAccount(login string) (*customer_account.CustomerAccount, error)
	GetCustomerAccounts() (customer_account.CustomerAccounts, error)

	AddOrder(userId string, record *message.AccuralMessage) error

	GetOrder(orderId string) (*order.Order, error)
	ListOrders(userId string) ([]order.Order, error)
	GetOrderOwner(orderId string) (string, error)
	UpdateOrder(order *order.Order) error

	AddOrderToPendingList(orderId string) error
	GetPendingOrders() ([]string, error)
	RemoveOrderFromPendingList(orderId string) error

	RegisterWithdrawn(userId string, withdraw *withdrawn.WithdrawnRecord) error
	GetWithdrawn(orderId string) (*withdrawn.WithdrawnRecord, error)
	ListWithdrawns(orderId string) ([]withdrawn.WithdrawnRecord, error)
}

func New(config *config.Config) Storage {
	var storage Storage
	logging.Debug("Initialization Storage")
	if len(config.DatabaseURI) == 0 {
		storage = memory.New()
	} else {
		storage = database.New()
	}
	logging.Debug("Storage Initialization finished")
	return storage
}
