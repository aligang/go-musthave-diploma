package memory

import (
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/orderrecord"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"sync"
)

type Storage struct {
	CustomerAccounts   account.CustomerAccounts
	Orders             orderrecord.Orders
	Withdrawns         withdrawn.Withdrawns
	PendingOrders      map[string]*struct{}
	CustomerOrders     map[string][]string
	CustomerWithdrawns map[string][]string
	Lock               sync.Mutex
}

func New() *Storage {
	logging.Debug("Initialization In-Memory Storage Backend")
	m := &Storage{
		CustomerAccounts:   account.CustomerAccounts{},
		Orders:             orderrecord.Orders{},
		Withdrawns:         withdrawn.Withdrawns{},
		PendingOrders:      map[string]*struct{}{},
		CustomerOrders:     map[string][]string{},
		CustomerWithdrawns: map[string][]string{},
	}
	logging.Debug("Initialization In-Memory Storage Backend is Finished")
	return m
}

func Init(
	customerAccounts account.CustomerAccounts,
	orders orderrecord.Orders,
	withdrawns withdrawn.Withdrawns,
	pendingOrders map[string]*struct{},
	customerOrders map[string][]string,
	customerWithdrawns map[string][]string,
) *Storage {
	m := &Storage{
		CustomerAccounts:   customerAccounts,
		Orders:             orders,
		Withdrawns:         withdrawns,
		PendingOrders:      pendingOrders,
		CustomerOrders:     customerOrders,
		CustomerWithdrawns: customerWithdrawns,
	}
	return m
}
