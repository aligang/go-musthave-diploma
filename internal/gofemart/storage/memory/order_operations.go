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

func (s *Storage) AddOrderToPendingList(ctx context.Context, orderId string) error {
	if _, exists := s.PendingOrders[orderId]; exists {
		return fmt.Errorf("record Already present")
	}
	s.PendingOrders[orderId] = nil
	return nil
}

func (s *Storage) GetPendingOrders(ctx context.Context) ([]string, error) {
	var pendingOrders []string
	for orderId, _ := range s.PendingOrders {
		pendingOrders = append(pendingOrders, orderId)
	}
	return pendingOrders, nil
}

func (s *Storage) RemoveOrderFromPendingList(ctx context.Context, orderId string) error {
	if _, exists := s.PendingOrders[orderId]; !exists {
		return fmt.Errorf("record does not exist")
	}
	delete(s.PendingOrders, orderId)
	return nil
}

func (s *Storage) GetOrder(orderId string) (*order.Order, error) {
	orderRecord, exists := s.Orders[orderId]
	if !exists {
		return nil, repository_errors.ErrNoContent
	}
	return orderRecord.Order, nil
}

func (s *Storage) GetOrderWithinTransaction(ctx context.Context, orderId string) (*order.Order, error) {
	return s.GetOrder(orderId)
}

func (s *Storage) ListOrders(userId string) ([]order.Order, error) {
	orderIds, exists := s.CustomerOrders[userId]
	if !exists {
		return []order.Order{}, nil
	}

	var err error
	var orders []order.Order
	for _, orderId := range orderIds {
		internalOrder, exists := s.Orders[orderId]
		if !exists {
			logging.Warn("order info for orderId=%s was not found, seems as DB data lost", orderId)
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
