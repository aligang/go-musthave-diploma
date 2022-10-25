package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) AddOrder(ctx context.Context, userID string, order *order.Order, tx *sqlx.Tx) error {
	s.log.Debug("Preparing statement to add order to Repository: %+v for user %s", order, userID)
	query := "INSERT INTO orders (Number, Status, Accural, UploadedAt, Owner) VALUES($1, $2, $3, $4, $5)"
	var args = []interface{}{order.Number, order.Status, order.Accural, order.UploadedAt, userID}
	return s.modifyOrder(ctx, order, query, args, tx)
}

func (s *Storage) UpdateOrder(ctx context.Context, order *order.Order, tx *sqlx.Tx) error {
	s.log.Debug("Preparing statement to update order to Repository: %+v", order)
	query := "UPDATE orders SET number = $1, status = $2, accural = $3, uploadedat = $4 WHERE number = $5"
	var args = []interface{}{order.Number, order.Status, order.Accural, order.UploadedAt, order.Number}
	return s.modifyOrder(ctx, order, query, args, tx)
}

func (s *Storage) modifyOrder(ctx context.Context, order *order.Order, query string, args []interface{}, tx *sqlx.Tx) error {

	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return err
	}
	s.log.Debug("Executing statement to modify order  Repository: %s %s", query, args)
	_, err = statement.ExecContext(ctx, args...)
	if err != nil {
		s.log.Warn("Error During statement Execution %s with %s, %s, %s, %s",
			query, args[0], args[1], args[2], args[3])
		return err
	}
	s.log.Debug("Order record update succseeded : %s", order.Number)
	return nil
}

func (s *Storage) GetOrder(ctx context.Context, orderID string, tx *sqlx.Tx) (*order.Order, error) {
	query := "SELECT number, status, accural, uploadedat FROM orders WHERE Number = $1"
	var args = []interface{}{orderID}
	s.log.Debug("Preparing statement to fetch order from Repository: %s", query)
	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return nil, err
	}
	s.log.Debug("Executing statement to fetch order info Repository: %s %s", query, orderID)
	orderInstance := &order.Order{}
	err = statement.GetContext(ctx, orderInstance, args...)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, repositoryerrors.ErrNoContent
	case err != nil:
		s.log.Warn("Error during decoding database response: %s", err.Error())
		return nil, err
	}

	return orderInstance, nil
}

func (s *Storage) ListOrders(ctx context.Context, userID string) ([]order.Order, error) {
	s.log.Debug("Preparing statement to fetch orders from Repository")
	query := "SELECT number, status, accural, uploadedat  FROM orders where owner = $1"
	args := []interface{}{userID}
	var orders []order.Order

	statement, err := s.DB.Preparex(query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return orders, err
	}
	s.log.Debug("Executing statement to fetch orders from Repository")
	err = statement.SelectContext(ctx, &orders, args...)
	if err != nil {
		s.log.Warn("Error During statement Execution %s with %s", query, args[0])
		return orders, err
	}
	return orders, nil
}

func (s *Storage) AddOrderToPendingList(ctx context.Context, orderID string, tx *sqlx.Tx) error {
	s.log.Debug("Preparing statement to delete pending order From Repository:  %s", orderID)
	query := "INSERT INTO pending_orders (order_id) VALUES($1)"
	var args = []interface{}{orderID}

	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return err
	}
	s.log.Debug("Executing statement to delete pending order from Repository: %+s", orderID)
	_, err = statement.ExecContext(ctx, args...)
	if err != nil {
		s.log.Warn("Error During statement Execution %s with %s", query, args[0])
		return err
	}
	return nil
}
func (s *Storage) RemoveOrderFromPendingList(ctx context.Context, orderID string, tx *sqlx.Tx) error {
	s.log.Debug("Preparing statement to delete pending order to Repository:  %s", orderID)
	query := "DELETE FROM pending_orders WHERE order_id = $1"
	var args = []interface{}{orderID}

	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return err
	}
	s.log.Debug("Executing statement to add pending order to Repository: %+s", orderID)
	_, err = statement.ExecContext(ctx, args...)
	if err != nil {
		s.log.Warn("Error During statement Execution %s with %s", query, args[0])
		return err
	}
	return nil
}

func (s *Storage) GetPendingOrders(ctx context.Context, tx *sqlx.Tx) ([]string, error) {
	s.log.Debug("Preparing statement to fetch pending order from Repository")
	query := "SELECT * FROM pending_orders"
	var args []interface{}
	var orders []string

	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return orders, err
	}
	s.log.Debug("Executing statement to fetch pending orders from Repository")
	err = statement.SelectContext(ctx, &orders, args...)

	if err != nil {
		s.log.Warn("Error During statement Execution %s with %s", query, err.Error())
		return orders, err
	}

	return orders, nil
}
