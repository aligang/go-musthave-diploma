package database

import "github.com/aligang/go-musthave-diploma/internal/withdrawn"

func (s *Storage) RegisterWithdrawn(userId string, withdraw *withdrawn.WithdrawnRecord) error {
	return nil
}

func (s *Storage) GetWithdrawn(orderId string) (*withdrawn.WithdrawnRecord, error) {
	return &withdrawn.WithdrawnRecord{}, nil
}

func (s *Storage) ListWithdrawns(userId string) ([]withdrawn.WithdrawnRecord, error) {
	return []withdrawn.WithdrawnRecord{}, nil
}
