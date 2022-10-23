package memory

import (
	"context"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
)

func (s *Storage) RegisterWithdrawn(ctx context.Context, userID string, withdrawn *withdrawn.WithdrawnRecord) error {
	s.Withdrawns[withdrawn.OrderId] = *withdrawn
	s.CustomerWithdrawns[userID] = append(s.CustomerWithdrawns[userID], withdrawn.OrderId)
	return nil
}

func (s *Storage) GetWithdrawnWithinTransaction(ctx context.Context, withdrawnID string) (*withdrawn.WithdrawnRecord, error) {
	withdrawn, exists := s.Withdrawns[withdrawnID]
	if !exists {
		return nil, repositoryerrors.ErrNoContent
	}
	return &withdrawn, nil
}

func (s *Storage) GetWithdrawnIds(userID string) ([]string, error) {
	return s.CustomerWithdrawns[userID], nil
}

func (s *Storage) ListWithdrawns(userID string) ([]withdrawn.WithdrawnRecord, error) {
	withdrawnIDs, exists := s.CustomerWithdrawns[userID]
	if !exists {
		return []withdrawn.WithdrawnRecord{}, nil
	}

	var err error
	var withdrawns []withdrawn.WithdrawnRecord
	for _, id := range withdrawnIDs {
		withdrawn, exists := s.Withdrawns[id]
		if !exists {
			logging.Warn("order info for orderID=%s was not found, seems as DB data lost", id)
			err = fmt.Errorf("porblem during fetching list of orders")
		}
		withdrawns = append(withdrawns, withdrawn)
	}
	return withdrawns, err
}
