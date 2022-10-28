package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) WithinTransaction(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin trx: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			s.log.Warn("rollback tx: %s", err.Error())
		}
		return fmt.Errorf("run tx: %w", err)
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			s.log.Warn("rollback tx: %s", err.Error())
		}
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
