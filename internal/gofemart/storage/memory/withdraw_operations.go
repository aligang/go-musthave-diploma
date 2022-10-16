package memory

import (
	"context"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
)

func (s *Storage) RegisterWithdrawn(ctx context.Context, userID string, withdrawn *withdrawn.WithdrawnRecord) error {
	s.Withdrawns[withdrawn.Order] = *withdrawn
	s.CustomerWithdrawns[userID] = append(s.CustomerWithdrawns[userID], withdrawn.Order)
	return nil
}

func (s *Storage) GetWithdrawnWithinTransaction(ctx context.Context, withdrawnId string) (*withdrawn.WithdrawnRecord, error) {
	withdrawn, exists := s.Withdrawns[withdrawnId]
	if !exists {
		return nil, repositoryerrors.ErrNoContent
	}
	return &withdrawn, nil
}

func (s *Storage) GetWithdrawnIds(userID string) ([]string, error) {
	return s.CustomerWithdrawns[userID], nil
}

func (s *Storage) ListWithdrawns(userID string) ([]withdrawn.WithdrawnRecord, error) {
	withdrawnIds, exists := s.CustomerWithdrawns[userID]
	if !exists {
		return []withdrawn.WithdrawnRecord{}, nil
	}

	var err error
	var withdrawns []withdrawn.WithdrawnRecord
	for _, id := range withdrawnIds {
		withdrawn, exists := s.Withdrawns[id]
		if !exists {
			logging.Warn("order info for orderID=%s was not found, seems as DB data lost", id)
			err = fmt.Errorf("porblem during fetching list of orders")
		}
		withdrawns = append(withdrawns, withdrawn)
	}
	return withdrawns, err
}
