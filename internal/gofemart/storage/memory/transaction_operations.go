package memory

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) WithinTransaction(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	s.Lock.Lock()
	err := fn(ctx, nil)
	if err != nil {
		return err
	}
	s.Lock.Unlock()
	return nil
}
