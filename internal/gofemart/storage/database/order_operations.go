package database

import (
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
)

func (s *Storage) AddOrder(userId string, record *message.AccuralMessage) error {
	return nil
}

func (s *Storage) GetOrder(orderId string) (*order.Order, error) {
	return nil, nil
}

func (s *Storage) ListOrders(userId string) ([]order.Order, error) {
	return []order.Order{}, nil
}

func (s *Storage) AddOrderToPendingList(orderId string) error {
	return nil
}
func (s *Storage) RemoveOrderFromPendingList(orderId string) error {
	return nil
}

func (s *Storage) UpdateOrder(order *order.Order) error {
	return nil
}

func (s *Storage) GetPendingOrders() ([]string, error) {
	return []string{}, nil
}
