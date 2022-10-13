package tracker

import (
	"github.com/aligang/go-musthave-diploma/internal/accural"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order/status"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"time"
)

type Tracker struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *Tracker {
	logging.Debug("Initializing pending list tracker")
	tracker := Tracker{
		storage: s,
		config:  cfg,
	}
	logging.Debug("Pending list tracker initialization succeeded")
	return &tracker
}

func (t *Tracker) RunInBackground() {
	go func() {
		ticker := time.NewTicker(TRACKING_INTERVALL)
		for {
			<-ticker.C
			t.Sync()
		}
	}()
}

func (t *Tracker) Sync() {
	var dbErr error
	logging.Warn("Checking state of orders in pending list")
	t.storage.StartTransaction()
	defer func() {
		if dbErr != nil {
			t.storage.RollbackTransaction()
			return
		}
		t.storage.CommitTransaction()
	}()

	orderIdsRecords, dbErr := t.storage.GetPendingOrders()
	if dbErr != nil {
		logging.Warn("Tracker failed to fetch list of pending orders from DB")
		return
	}
	if len(orderIdsRecords) == 0 {
		logging.Debug("Tracker has nothing to check: Pending list is empty")
		return
	}

	orderIds := make([]string, len(orderIdsRecords))
	var updatedOrdersCounter uint64
	var proceededOrdersCounter uint64
	copy(orderIds, orderIdsRecords)
	logging.Debug("There %d record(s) in pending list", len(orderIds))

	for _, orderId := range orderIds {
		accuralRecord, err := accural.FetchOrderInfo(orderId, t.config)
		if err != nil {
			logging.Warn("Tracker failed to fetch accural info for order %s", orderId)
			continue
		}
		order, dbErr := t.storage.GetOrderWithinTransaction(orderId)
		if dbErr != nil {
			logging.Warn("Tracker failed to fetch order info from DB %s", orderId)
			return
		}
		if order.Status != accuralRecord.Status {
			logging.Warn("Order %s status needs to be updated", orderId)
			order.Status = accuralRecord.Status
			dbErr = t.storage.UpdateOrder(order)
			if dbErr != nil {
				logging.Warn("Could not update order record %s", orderId)
				return
			}
			updatedOrdersCounter += 1
		}
		if order.Status == status.PROCESSED {
			userId, dbErr := t.storage.GetOrderOwner(orderId)
			if dbErr != nil {
				logging.Warn("Tracker could not fetch owner of order: %s", orderId)
				return
			}
			accountInfo, dbErr := t.storage.GetCustomerAccount(userId)
			if dbErr != nil {
				logging.Warn("Tracker could not fetch account info of : %s", userId)
				return
			}
			dbErr = t.storage.RemoveOrderFromPendingList(orderId)
			if dbErr != nil {
				logging.Warn("Tracker failed to remove order from pending list : %s", orderId)
				return
			}
			accountInfo.Current += order.Accural
			dbErr = t.storage.UpdateCustomerAccount(accountInfo)
			if dbErr != nil {
				logging.Warn("Tracker failed to update account info : %s", userId)
				return
			}
			proceededOrdersCounter += 1
		}
	}
	logging.Warn("Pending list check succeeded")
	if updatedOrdersCounter == 0 {
		logging.Warn("no order records were updated")
		return
	}
	logging.Warn("%d order record(s) were updated", updatedOrdersCounter)
	logging.Warn("%d  new order record(s) are now in PROCEEDED status", proceededOrdersCounter)
}
