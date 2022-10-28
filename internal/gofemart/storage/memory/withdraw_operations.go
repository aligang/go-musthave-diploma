package memory

import (
	"context"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) RegisterWithdrawn(ctx context.Context, userID string, withdrawn *withdrawn.WithdrawnRecord, tx *sqlx.Tx) error {
	s.Withdrawns[withdrawn.OrderID] = *withdrawn
	s.CustomerWithdrawns[userID] = append(s.CustomerWithdrawns[userID], withdrawn.OrderID)
	return nil
}

func (s *Storage) GetWithdrawn(ctx context.Context, withdrawnID string, tx *sqlx.Tx) (*withdrawn.WithdrawnRecord, error) {
	withdrawn, exists := s.Withdrawns[withdrawnID]
	if !exists {
		return nil, repositoryerrors.ErrNoContent
	}
	return &withdrawn, nil
}

func (s *Storage) GetWithdrawnIds(ctx context.Context, userID string) ([]string, error) {
	return s.CustomerWithdrawns[userID], nil
}

func (s *Storage) ListWithdrawns(ctx context.Context, userID string) ([]withdrawn.WithdrawnRecord, error) {
	withdrawnIDs, exists := s.CustomerWithdrawns[userID]
	if !exists {
		return []withdrawn.WithdrawnRecord{}, nil
	}

	var err error
	var withdrawns []withdrawn.WithdrawnRecord
	for _, id := range withdrawnIDs {
		withdrawn, exists := s.Withdrawns[id]
		if !exists {
			s.log.Warn("order info for orderID=%s was not found, seems as DB data lost", id)
			err = fmt.Errorf("porblem during fetching list of orders")
		}
		withdrawns = append(withdrawns, withdrawn)
	}
	return withdrawns, err
}
