package store

import (
	"errors"
	"sync"
)

// Storage defines the requirements for a type which can persist and retrieve key/value pairs.
type Storage interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

// Store is a concurrency safe map-driven key/value store. It satisfies the Storage interface.
type Store struct {
	lookup map[string]string
	mu     *sync.RWMutex
}

// Ensure Store satisfies Storage.
var _ Storage = Store{}

// New creates an initialised Store.
func New() Store {
	return Store{
		lookup: make(map[string]string),
		mu:     &sync.RWMutex{},
	}
}

// Set sets the given value for a given key in the store. If the key exists already, the value will be overwritten.
func (s Store) Set(key, value string) error {
	s.mu.Lock()
	s.lookup[key] = value
	s.mu.Unlock()
	return nil
}

// ErrKeyNotFound indicates that a value could not be found in the store for the provided key.
var ErrKeyNotFound = errors.New("key not found in store")

// Get returns the value for a given key. If the key is not found, ErrKeyNotFound is returned.
func (s Store) Get(key string) (string, error) {
	s.mu.RLock()
	val, ok := s.lookup[key]
	s.mu.RUnlock()

	if !ok {
		return "", ErrKeyNotFound
	}
	return val, nil
}
