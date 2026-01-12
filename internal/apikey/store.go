package apikey

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type Store struct {
	mu sync.RWMutex
	keys map[string]struct{}
}

func NewStore() *Store{
	return &Store{keys: make(map[string]struct{})}
}

func (s *Store) Create() (string, error) {
	// 32 bytes => 64 hex chars
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	key := hex.EncodeToString(buf)

	s.mu.Lock()
	s.keys[key] = struct{}{}
	s.mu.Unlock()

	return key, nil
}

func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	_, ok := s.keys[key]
	s.mu.RUnlock()
	return ok
}

func (s *Store) Revoke(key string) bool {
	s.mu.Lock()
	_, ok := s.keys[key]
	if ok {
		delete(s.keys, key)
	}
	s.mu.Unlock()
	return ok
}