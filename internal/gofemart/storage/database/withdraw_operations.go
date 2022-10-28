package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) RegisterWithdrawn(ctx context.Context, userID string, withdraw *withdrawn.WithdrawnRecord, tx *sqlx.Tx) error {
	s.log.Debug("Preparing statement to add withdraw to Repository: %+v for user %s", withdraw, userID)
	query := "INSERT INTO withdrawns (OrderID, Sum, ProcessedAt, Owner) VALUES($1, $2, $3, $4)"
	var args = []interface{}{withdraw.OrderID, withdraw.Sum, withdraw.ProcessedAt, userID}
	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return err
	}
	s.log.Debug("Executing statement to modify order  Repository: %s with %s, %s, %s, %s ",
		query, args[0], args[1], args[2], args[3])
	_, err = statement.ExecContext(ctx, args...)
	if err != nil {
		s.log.Warn("Error During statement Execution %s with %d, %f, %s,",
			query, args[0], args[1], args[2])
		return err
	}
	s.log.Debug("Order Withdraw add succeeded : %+v", withdraw)
	return nil

}

func (s *Storage) GetWithdrawn(ctx context.Context, orderID string, tx *sqlx.Tx) (*withdrawn.WithdrawnRecord, error) {
	query := "SELECT OrderID, Sum, ProcessedAt FROM withdrawns WHERE OrderID = $1"
	var args = []interface{}{orderID}
	s.log.Debug("Preparing statement to fetch order from Repository: %s", query)
	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return nil, err
	}
	s.log.Debug("Executing statement to fetch withdraw info from Repository: %s %s", query, orderID)
	withdrawnInstance := &withdrawn.WithdrawnRecord{
		Withdrawn: &withdrawn.Withdrawn{},
	}
	err = statement.GetContext(ctx, withdrawnInstance, args...)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		s.log.Warn("Database response is empty")
		return nil, repositoryerrors.ErrNoContent
	case err != nil:
		s.log.Warn("Error During statement Execution %s with %s", query, args[0])
		return nil, err
	default:
		return withdrawnInstance, nil
	}
}

func (s *Storage) ListWithdrawns(ctx context.Context, userID string) ([]withdrawn.WithdrawnRecord, error) {
	s.log.Debug("Preparing statement to fetch orders from Repository")
	query := "SELECT OrderID, Sum, ProcessedAt FROM withdrawns WHERE Owner = $1"
	args := []interface{}{userID}
	var withdrawns []withdrawn.WithdrawnRecord

	statement, err := s.DB.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return withdrawns, err
	}
	s.log.Debug("Executing statement to fetch withdrawns from Repository")

	err = statement.SelectContext(ctx, &withdrawns, args...)
	if err != nil {
		s.log.Warn("Error During statement Execution %s with %s", query, args[0])
		return withdrawns, err
	}
	return withdrawns, nil
}
