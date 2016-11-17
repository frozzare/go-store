package rwmutex

import (
	"sync"

	"github.com/frozzare/go-store/driver"
)

var (
	instancesMu sync.RWMutex
	instances   = make(map[string]driver.Driver)
)

// Driver represents a rwmutex driver.
type Driver struct {
	lock sync.RWMutex
	data map[string][]byte
}

// Open creates a new RWMutex store.
func Open(args ...interface{}) driver.Driver {
	name := ""
	if len(args) > 0 {
		name = args[0].(string)
	}

	instancesMu.Lock()

	defer instancesMu.Unlock()

	if instances[name] == nil {
		instances[name] = &Driver{data: make(map[string][]byte)}
	}

	return instances[name]
}

// Open creates a new RWMutex store with a specified instance.
func (s *Driver) Open(args ...interface{}) driver.Driver {
	return Open(args...)
}

// Count returns numbers of keys in store.
func (s *Driver) Count() int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return int64(len(s.data))
}

// Exists returns true when a key exists false when not existing in store.
func (s *Driver) Exists(key string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, exists := s.data[key]
	return exists
}

// Get value from key in store.
func (s *Driver) Get(key string) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.data[key], nil
}

// Set key with value in store.
func (s *Driver) Set(key string, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
	return nil
}

// Delete key from store.
func (s *Driver) Delete(key string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
	return nil
}

// Close does not exists for RWMutex driver.
func (s *Driver) Close() error {
	return nil
}
