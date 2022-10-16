package memory

import (
	"context"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/internal_order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repository_errors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
)

func (s *Storage) AddOrder(ctx context.Context, userId string, order *order.Order) error {
	orderInstance := internal_order.New(userId, order)
	logging.Debug("order to be Stored: %+v", orderInstance)
	s.Orders[order.Number] = orderInstance
	s.CustomerOrders[userId] = append(s.CustomerOrders[userId], order.Number)
	return nil
}

func (s *Storage) AddOrderToPendingList(ctx context.Context, orderID string) error {
	if _, exists := s.PendingOrders[orderID]; exists {
		return fmt.Errorf("record Already present")
	}
	s.PendingOrders[orderID] = nil
	return nil
}

func (s *Storage) GetPendingOrders(ctx context.Context) ([]string, error) {
	var pendingOrders []string
	for orderID, _ := range s.PendingOrders {
		pendingOrders = append(pendingOrders, orderID)
	}
	return pendingOrders, nil
}

func (s *Storage) RemoveOrderFromPendingList(ctx context.Context, orderID string) error {
	if _, exists := s.PendingOrders[orderID]; !exists {
		return fmt.Errorf("record does not exist")
	}
	delete(s.PendingOrders, orderID)
	return nil
}

func (s *Storage) GetOrder(orderID string) (*order.Order, error) {
	orderRecord, exists := s.Orders[orderID]
	if !exists {
		return nil, repository_errors.ErrNoContent
	}
	return orderRecord.Order, nil
}

func (s *Storage) GetOrderWithinTransaction(ctx context.Context, orderID string) (*order.Order, error) {
	return s.GetOrder(orderID)
}

func (s *Storage) ListOrders(userId string) ([]order.Order, error) {
	orderIDs, exists := s.CustomerOrders[userId]
	if !exists {
		return []order.Order{}, nil
	}

	var err error
	var orders []order.Order
	for _, orderID := range orderIDs {
		internalOrder, exists := s.Orders[orderID]
		if !exists {
			logging.Warn("order info for orderID=%s was not found, seems as DB data lost", orderID)
			err = fmt.Errorf("porblem during fetching list of orders")
		}
		orders = append(orders, *internalOrder.Order)
	}
	return orders, err
}

func (s *Storage) UpdateOrder(ctx context.Context, order *order.Order) error {
	orderInstance, exists := s.Orders[order.Number]
	if !exists {
		return fmt.Errorf("order record was not found in database")
	}
	orderInstance.Status = order.Status
	s.Orders[order.Number] = orderInstance
	return nil
}
