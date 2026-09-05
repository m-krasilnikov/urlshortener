package storage

import "sync"

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) Save(id, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = url
}

func (s *MemoryStorage) Get(id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.data[id]

	return url, ok
}
