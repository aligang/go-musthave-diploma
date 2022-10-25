package tracker

import (
	"context"
	"github.com/aligang/go-musthave-diploma/internal/accural"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order/status"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/jmoiron/sqlx"
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
	logging.Warn("Checking state of orders in pending list")
	var updatedOrdersCounter uint64
	var proceededOrdersCounter uint64

	t.storage.WithinTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		orderIDsRecords, err := t.storage.GetPendingOrders(ctx, tx)
		if err != nil {
			logging.Warn("Tracker failed to fetch list of pending orders from DB")
			return err
		}
		if len(orderIDsRecords) == 0 {
			logging.Debug("Tracker has nothing to check: Pending list is empty")
			return nil
		}

		orderIDs := make([]string, len(orderIDsRecords))

		copy(orderIDs, orderIDsRecords)
		logging.Debug("There %d record(s) in pending list", len(orderIDs))

		for _, orderID := range orderIDs {
			accuralRecord, err := t.accrualClient.FetchOrderInfo(ctx, orderID)
			if err != nil {
				logging.Warn("Tracker failed to fetch accural info for order %s", orderID)
				continue
			}
			order, err := t.storage.GetOrder(ctx, orderID, tx)
			if err != nil {
				logging.Warn("Tracker failed to fetch order info from DB %s", orderID)
				return err
			}
			if order.Status != accuralRecord.Status {
				logging.Warn("Order %s status needs to be updated", orderID)
				order.Status = accuralRecord.Status
				err = t.storage.UpdateOrder(ctx, order, tx)
				if err != nil {
					logging.Warn("Could not update order record %s", orderID)
					return err
				}
				updatedOrdersCounter += 1
			}
			if order.Status == status.PROCESSED || order.Status == status.INVALID {
				userID, err := t.storage.GetOrderOwner(ctx, orderID, tx)
				if err != nil {
					logging.Warn("Tracker could not fetch owner of order: %s", orderID)
					return err
				}
				accountInfo, err := t.storage.GetCustomerAccount(ctx, userID, tx)
				if err != nil {
					logging.Warn("Tracker could not fetch account info of : %s", userID)
					return err
				}
				err = t.storage.RemoveOrderFromPendingList(ctx, orderID, tx)
				if err != nil {
					logging.Warn("Tracker failed to remove order from pending list : %s", orderID)
					return err
				}
				if order.Status == status.PROCESSED {
					accountInfo.Current += order.Accural
					err = t.storage.UpdateCustomerAccount(ctx, accountInfo, tx)
					if err != nil {
						logging.Warn("Tracker failed to update account info : %s", userID)
						return err
					}
				}
				proceededOrdersCounter += 1
			}
		}
		logging.Warn("Pending list check succeeded")
		if updatedOrdersCounter == 0 {
			logging.Warn("no order records were updated")
			return nil
		}
		logging.Warn("%d order record(s) were updated", updatedOrdersCounter)
		logging.Warn("%d  new order record(s) are now in PROCEEDED status", proceededOrdersCounter)
		return nil
	})
}
