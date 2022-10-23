package tracker

import (
	"context"
	"github.com/aligang/go-musthave-diploma/internal/accural"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order/status"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"sync"
	"time"
)

type Tracker struct {
	storage       storage.Storage
	config        *config.Config
	accrualClient *accural.AccrualClient
}

func New(s storage.Storage, cfg *config.Config) *Tracker {
	logging.Debug("Initializing pending list tracker")
	tracker := Tracker{
		storage:       s,
		config:        cfg,
		accrualClient: accural.New(cfg),
	}
	logging.Debug("Pending list tracker initialization succeeded")
	return &tracker
}

func (t *Tracker) RunPeriodically(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	ticker := time.NewTicker(trackingInterval)
	for {
		<-ticker.C
		t.Sync(ctx)
		<-ctx.Done()
		wg.Add(-1)
	}
}

func (t *Tracker) Sync(ctx context.Context) {
	var dbErr error
	logging.Warn("Checking state of orders in pending list")

	t.storage.StartTransaction(ctx)
	defer func() {
		if dbErr != nil {
			t.storage.RollbackTransaction(ctx)
			return
		}
		t.storage.CommitTransaction(ctx)
	}()

	select {
	default:
	case <-ctx.Done():
		return
	}
	orderIDsRecords, dbErr := t.storage.GetPendingOrders(ctx)
	if dbErr != nil {
		logging.Warn("Tracker failed to fetch list of pending orders from DB")
		return
	}
	if len(orderIDsRecords) == 0 {
		logging.Debug("Tracker has nothing to check: Pending list is empty")
		return
	}

	orderIDs := make([]string, len(orderIDsRecords))
	var updatedOrdersCounter uint64
	var proceededOrdersCounter uint64
	copy(orderIDs, orderIDsRecords)
	logging.Debug("There %d record(s) in pending list", len(orderIDs))

	for _, orderID := range orderIDs {
		select {
		default:
		case <-ctx.Done():
			return
		}
		accuralRecord, err := t.accrualClient.FetchOrderInfo(ctx, orderID)
		if err != nil {
			logging.Warn("Tracker failed to fetch accural info for order %s", orderID)
			continue
		}
		order, dbErr := t.storage.GetOrderWithinTransaction(ctx, orderID)
		if dbErr != nil {
			logging.Warn("Tracker failed to fetch order info from DB %s", orderID)
			return
		}
		if order.Status != accuralRecord.Status {
			logging.Warn("Order %s status needs to be updated", orderID)
			order.Status = accuralRecord.Status
			dbErr = t.storage.UpdateOrder(ctx, order)
			if dbErr != nil {
				logging.Warn("Could not update order record %s", orderID)
				return
			}
			updatedOrdersCounter += 1
		}
		if order.Status == status.PROCESSED || order.Status == status.INVALID {
			userID, dbErr := t.storage.GetOrderOwner(ctx, orderID)
			if dbErr != nil {
				logging.Warn("Tracker could not fetch owner of order: %s", orderID)
				return
			}
			accountInfo, dbErr := t.storage.GetCustomerAccount(userID)
			if dbErr != nil {
				logging.Warn("Tracker could not fetch account info of : %s", userID)
				return
			}
			dbErr = t.storage.RemoveOrderFromPendingList(ctx, orderID)
			if dbErr != nil {
				logging.Warn("Tracker failed to remove order from pending list : %s", orderID)
				return
			}
			if order.Status == status.PROCESSED {
				accountInfo.Current += order.Accural
				dbErr = t.storage.UpdateCustomerAccount(ctx, accountInfo)
				if dbErr != nil {
					logging.Warn("Tracker failed to update account info : %s", userID)
					return
				}
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
