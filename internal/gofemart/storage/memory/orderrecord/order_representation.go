package orderrecord

import (
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
)

type Order struct {
	*order.Order
	Owner string
}

type Orders map[string]Order

func New(userID string, order *order.Order) Order {
	return Order{
		Order: order,
		Owner: userID,
	}
}
