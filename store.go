package store

import "sync"

var instances map[string]*Store

func init() {
	instances = make(map[string]*Store)
}

// Store represents the store.
type Store struct {
	lock sync.RWMutex
	data map[string]interface{}
}

// New creates a new store.
func New() *Store {
	return &Store{data: make(map[string]interface{})}
}

// Instance return store instance.
func Instance(params ...string) *Store {
	key := ""

	if len(params) > 0 {
		key = params[0]
	}

	if instances[key] == nil {
		instances[key] = New()
	}

	return instances[key]
}

// Count returns numbers of keys in store.
func (s *Store) Count() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.data)
}

// Exists returns true when a key exists false when not existing in store.
func (s *Store) Exists(key string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, exists := s.data[key]
	return exists
}

// Get value from key in store.
func (s *Store) Get(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.data[key]
}

// Set key with value in store.
func (s *Store) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

// Delete key from store.
func (s *Store) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
}
