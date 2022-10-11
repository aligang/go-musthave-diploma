package memory

func (s *Storage) StartTransaction() {
	s.Lock.Lock()
}

func (s *Storage) CommitTransaction() {
	s.Lock.Unlock()
}

func (s *Storage) RollbackTransaction() {
	s.Lock.Unlock()
}
