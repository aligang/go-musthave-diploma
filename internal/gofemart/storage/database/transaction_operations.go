package database

import (
	"github.com/aligang/go-musthave-diploma/internal/logging"
)

func (s *Storage) StartTransaction() {
	tx, err := s.DB.Begin()
	if err != nil {
		logging.Warn("Error during transaction creation: %s", err.Error())
		//return err
	}
	s.Tx = tx
}

func (s *Storage) CommitTransaction() {
	err := s.Tx.Commit()
	if err != nil {
		logging.Warn("Error during transaction commit: %s", err.Error())
	}
}

func (s *Storage) RollbackTransaction() {
	err := s.Tx.Rollback()
	if err != nil {
		logging.Warn("Error during transaction rollback: %s", err.Error())

	}
}
