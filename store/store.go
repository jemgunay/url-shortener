package store

import (
	"errors"
	"sync"
)

type Storage interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

type Store struct {
	lookup map[string]string
	mu     *sync.RWMutex
}

func New() Store {
	return Store{
		lookup: make(map[string]string),
		mu:     &sync.RWMutex{},
	}
}

func (s Store) Set(key, value string) error {
	s.mu.Lock()
	s.lookup[key] = value
	s.mu.Unlock()
	return nil
}

var ErrKeyNotFound = errors.New("key not found in store")

func (s Store) Get(key string) (string, error) {
	s.mu.RLock()
	val, ok := s.lookup[key]
	s.mu.RUnlock()

	if !ok {
		return "", ErrKeyNotFound
	}
	return val, nil
}
