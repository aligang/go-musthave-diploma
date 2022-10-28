package memory

import (
	"context"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/orderrecord"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) AddOrder(ctx context.Context, userID string, order *order.Order, tx *sqlx.Tx) error {
	orderInstance := orderrecord.New(userID, order)
	s.log.Debug("order to be Stored: %+v", orderInstance)
	s.Orders[order.Number] = orderInstance
	s.CustomerOrders[userID] = append(s.CustomerOrders[userID], order.Number)
	return nil
}

func (s *Storage) AddOrderToPendingList(ctx context.Context, orderID string, tx *sqlx.Tx) error {
	if _, exists := s.PendingOrders[orderID]; exists {
		return fmt.Errorf("record Already present")
	}
	s.PendingOrders[orderID] = nil
	return nil
}

func (s *Storage) GetPendingOrders(ctx context.Context, tx *sqlx.Tx) ([]string, error) {
	var pendingOrders []string
	for orderID := range s.PendingOrders {
		pendingOrders = append(pendingOrders, orderID)
	}
	return pendingOrders, nil
}

func (s *Storage) RemoveOrderFromPendingList(ctx context.Context, orderID string, tx *sqlx.Tx) error {
	if _, exists := s.PendingOrders[orderID]; !exists {
		return fmt.Errorf("record does not exist")
	}
	delete(s.PendingOrders, orderID)
	return nil
}

func (s *Storage) GetOrder(ctx context.Context, orderID string, tx *sqlx.Tx) (*order.Order, error) {
	orderRecord, exists := s.Orders[orderID]
	if !exists {
		return nil, repositoryerrors.ErrNoContent
	}
	return orderRecord.Order, nil
}

func (s *Storage) ListOrders(ctx context.Context, userID string) ([]order.Order, error) {
	orderIDs, exists := s.CustomerOrders[userID]
	if !exists {
		return []order.Order{}, nil
	}

	var err error
	var orders []order.Order
	for _, orderID := range orderIDs {
		internalOrder, exists := s.Orders[orderID]
		if !exists {
			s.log.Warn("order info for orderID=%s was not found, seems as DB data lost", orderID)
			err = fmt.Errorf("porblem during fetching list of orders")
		}
		orders = append(orders, *internalOrder.Order)
	}
	return orders, err
}

func (s *Storage) UpdateOrder(ctx context.Context, order *order.Order, tx *sqlx.Tx) error {
	orderInstance, exists := s.Orders[order.Number]
	if !exists {
		return fmt.Errorf("order record was not found in database")
	}
	orderInstance.Status = order.Status
	s.Orders[order.Number] = orderInstance
	return nil
}
