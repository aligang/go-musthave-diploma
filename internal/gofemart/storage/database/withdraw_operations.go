package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
)

func (s *Storage) RegisterWithdrawn(ctx context.Context, userID string, withdraw *withdrawn.WithdrawnRecord) error {
	logging.Debug("Preparing statement to add withdraw to Repository: %+v for user %s", withdraw, userID)
	query := "INSERT INTO withdrawns (Order_Id, Sum, Processed_at, Owner) VALUES($1, $2, $3, $4)"
	var args = []interface{}{withdraw.Order, withdraw.Sum, withdraw.ProcessedAt, userID}
	statement, err := s.Tx[ctx].Preparex(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return err
	}
	logging.Debug("Executing statement to modify order  Repository: %s with %s, %s, %s, %s ",
		query, args[0], args[1], args[2], args[3])
	_, err = statement.Exec(args...)
	if err != nil {
		logging.Warn("Error During statement Execution %s with %d, %f, %s,",
			query, args[0], args[1], args[2])
		return err
	}
	logging.Debug("Order Withdraw add succeeded : %+v", withdraw)
	return nil

}

func (s *Storage) GetWithdrawnWithinTransaction(ctx context.Context, orderID string) (*withdrawn.WithdrawnRecord, error) {
	query := "SELECT Order_Id, Sum, Processed_at FROM withdrawns WHERE Order_Id = $1"
	var args = []interface{}{orderID}
	logging.Debug("Preparing statement to fetch order from Repository: %s", query)
	statement, err := s.Tx[ctx].Preparex(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return nil, err
	}
	logging.Debug("Executing statement to fetch withdraw info from Repository: %s %s", query, orderID)
	withdrawnInstance := &withdrawn.WithdrawnRecord{
		Withdrawn: &withdrawn.Withdrawn{},
	}
	err = statement.Select(withdrawnInstance, args...)
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s", query, args[0])
		return nil, err
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		logging.Warn("Database response is empty")
		return nil, repositoryerrors.ErrNoContent
	case err != nil:
		logging.Warn("Error during decoding database response")
		return nil, err
	default:
		return withdrawnInstance, nil
	}
}

func (s *Storage) ListWithdrawns(userID string) ([]withdrawn.WithdrawnRecord, error) {
	logging.Debug("Preparing statement to fetch orders from Repository")
	query := "SELECT Order_Id, Sum, Processed_at FROM withdrawns WHERE Owner = $1"
	args := []interface{}{userID}
	var withdrawns []withdrawn.WithdrawnRecord

	statement, err := s.DB.Preparex(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return withdrawns, err
	}
	logging.Debug("Executing statement to fetch withdrawns from Repository")

	err = statement.Select(&withdrawns, args...)
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s", query, args[0])
		return withdrawns, err
	}
	return withdrawns, nil
}
