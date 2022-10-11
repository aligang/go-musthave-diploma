package memory

import (
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"sync"
)
import "github.com/aligang/go-musthave-diploma/internal/accural/message"

type Storage struct {
	AccuralMessageMap message.AccuralMessageMap
	Lock              sync.Mutex
}

func New() *Storage {
	logging.Debug("Initialization Storage")
	return &Storage{AccuralMessageMap: message.AccuralMessageMap{}}
}

func (s *Storage) Get(orderId string) (*message.AccuralMessage, error) {
	res, ok := s.AccuralMessageMap[orderId]
	if !ok {
		return nil, fmt.Errorf("record was not found")
	}
	return &res, nil
}

func (s *Storage) BulkGet() map[string]message.AccuralMessage {
	return s.AccuralMessageMap
}

func (s *Storage) Put(order *message.AccuralMessage) error {
	s.AccuralMessageMap[order.Order] = *order
	return nil
}

func Init(accuralMessages message.AccuralMessageMap) *Storage {
	s := &Storage{
		AccuralMessageMap: accuralMessages,
	}
	return s
}
