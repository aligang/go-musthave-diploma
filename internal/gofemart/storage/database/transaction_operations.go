package database

import (
	"context"
	"database/sql"
	"github.com/aligang/go-musthave-diploma/internal/logging"
)

func (s *Storage) StartTransaction(ctx context.Context) {
	tx, err := s.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		logging.Warn("Error during transaction creation: %s", err.Error())
		//return err
	}
	s.Lock.Lock()
	s.Tx[ctx] = tx
	s.Lock.Unlock()
}

func (s *Storage) CommitTransaction(ctx context.Context) {
	select {
	default:
		err := s.Tx[ctx].Commit()
		if err != nil {
			logging.Warn("Error during transaction commit: %s", err.Error())
		}
	case <-ctx.Done():
	}
	s.Lock.Lock()
	delete(s.Tx, ctx)
	s.Lock.Unlock()
}

func (s *Storage) RollbackTransaction(ctx context.Context) {
	select {
	default:
		err := s.Tx[ctx].Rollback()
		if err != nil {
			logging.Warn("Error during transaction rollback: %s", err.Error())

		}
	case <-ctx.Done():
	}
	s.Lock.Lock()
	delete(s.Tx, ctx)
	s.Lock.Unlock()
}
