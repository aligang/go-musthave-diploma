package database

import (
	"database/sql"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repository_errors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
)

func (s *Storage) AddOrder(userId string, order *order.Order) error {
	logging.Debug("Preparing statement to add order to Repository: %+v for user %s", order, userId)
	query := "INSERT INTO orders (Number, Status, Accural, UploadedAt, Owner) VALUES($1, $2, $3, $4, $5)"
	var args = []any{order.Number, order.Status, order.Accural, order.UploadedAt, userId}
	return s.modifyOrder(order, query, args)
}

func (s *Storage) UpdateOrder(order *order.Order) error {
	logging.Debug("Preparing statement to update order to Repository: %+v", order)
	query := "UPDATE orders SET number = $1, status = $2, accural = $3, uploadedat = $4 WHERE number = $5"
	var args = []any{order.Number, order.Status, order.Accural, order.UploadedAt, order.Number}
	return s.modifyOrder(order, query, args)
}

func (s *Storage) modifyOrder(order *order.Order, query string, args []any) error {

	statement, err := s.Tx.Prepare(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return err
	}
	logging.Debug("Executing statement to modify order  Repository: %s %s", query, args)
	_, err = statement.Exec(args...)
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s, %s, %s, $s",
			query, args[0], args[1], args[2], args[3])
		return err
	}
	logging.Debug("Order record update succseeded : %+s", order)
	return nil
}

func (s *Storage) GetOrder(orderId string) (*order.Order, error) {
	return s.getOrderCommon(orderId, s.DB.Prepare)
}

func (s *Storage) GetOrderWithinTransaction(orderId string) (*order.Order, error) {
	return s.getOrderCommon(orderId, s.Tx.Prepare)
}

func (s *Storage) getOrderCommon(orderId string, prepareFunc func(query string) (*sql.Stmt, error)) (*order.Order, error) {
	query := "SELECT number, status, accural, uploadedat FROM orders WHERE Number = $1"
	var args = []any{orderId}
	logging.Debug("Preparing statement to fetch order from Repository: %s", query)
	statement, err := prepareFunc(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return nil, err
	}
	logging.Debug("Executing statement to fetch order info Repository: %s %s", query, orderId)
	row := statement.QueryRow(args...)

	if row.Err() != nil {
		logging.Warn("Error During statement Execution %s with %s: %s", query, orderId, row.Err().Error())
		return nil, row.Err()
	}

	orderInstance := &order.Order{}
	err = row.Scan(&orderInstance.Number, &orderInstance.Status, &orderInstance.Accural, &orderInstance.UploadedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		logging.Warn("Database response is empty")
		return nil, repository_errors.ErrNoContent
	case err != nil:
		logging.Warn("Error during decoding database response")
		return nil, err
	default:
		return orderInstance, nil
	}
}

func (s *Storage) ListOrders(userId string) ([]order.Order, error) {
	logging.Debug("Preparing statement to fetch orders from Repository")
	query := "SELECT number, status, accural, uploadedat  FROM orders where owner = $1"
	args := []any{userId}
	var orders []order.Order

	statement, err := s.DB.Prepare(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return orders, err
	}
	logging.Debug("Executing statement to fetch orders from Repository")
	rows, err := statement.Query(args...)
	defer rows.Close()
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s", query, args[0])
		return orders, err
	}
	if err = rows.Err(); err != nil {
		logging.Warn("No records were returned from database")
		return orders, err
	}
	for rows.Next() {
		var orderInstance order.Order
		err = rows.Scan(&orderInstance.Number, &orderInstance.Status, &orderInstance.Accural, &orderInstance.UploadedAt)
		if err != nil {
			logging.Warn("problem during parsing data from repository")
			return orders, err
		}
		orders = append(orders, orderInstance)
	}

	return orders, nil
}

func (s *Storage) AddOrderToPendingList(orderId string) error {
	logging.Debug("Preparing statement to delete pending order From Repository:  %s", orderId)
	query := "INSERT INTO pending_orders (order_id) VALUES($1)"
	var args = []any{orderId}

	statement, err := s.Tx.Prepare(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return err
	}
	logging.Debug("Executing statement to delete pending order from Repository: %+s", orderId)
	_, err = statement.Exec(args...)
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s", query, args[0])
		return err
	}
	return nil
}
func (s *Storage) RemoveOrderFromPendingList(orderId string) error {
	logging.Debug("Preparing statement to delete pending order to Repository:  %s", orderId)
	query := "DELETE FROM pending_orders WHERE order_id = $1"
	var args = []any{orderId}

	statement, err := s.Tx.Prepare(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return err
	}
	logging.Debug("Executing statement to add pending order to Repository: %+s", orderId)
	_, err = statement.Exec(args...)
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s", query, args[0])
		return err
	}
	return nil
}

func (s *Storage) GetPendingOrders() ([]string, error) {
	logging.Debug("Preparing statement to fetch pending order from Repository")
	query := "SELECT * FROM pending_orders"
	var args []any
	var orders []string

	statement, err := s.Tx.Prepare(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return orders, err
	}
	logging.Debug("Executing statement to fetch pending orders from Repository")
	rows, err := statement.Query(args...)
	defer rows.Close()
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s", query, args[0])
		return orders, err
	}
	if err = rows.Err(); err != nil {
		logging.Warn("No records were returned from database")
		return orders, err
	}
	for rows.Next() {
		var orderId string
		err = rows.Scan(&orderId)
		if err != nil {
			logging.Warn("problem during parsing data from repository")
			return orders, err
		}
		orders = append(orders, orderId)
	}

	return orders, nil
}
