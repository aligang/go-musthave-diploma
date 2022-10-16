package memory

import "context"

func (s *Storage) StartTransaction(ctx context.Context) {
	s.Lock.Lock()

}

func (s *Storage) CommitTransaction(ctx context.Context) {
	s.Lock.Unlock()
}

func (s *Storage) RollbackTransaction(ctx context.Context) {
	s.Lock.Unlock()
}
