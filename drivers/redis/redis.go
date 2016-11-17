package redis

import (
	"github.com/frozzare/go-store/driver"
	"gopkg.in/redis.v3"
)

// Driver represents a Redis driver.
type Driver struct {
	client *redis.Client
}

// Open creates a new Redis store.
func Open(args ...interface{}) driver.Driver {
	var client *redis.Client

	if len(args) > 0 {
		client = args[0].(*redis.Client)
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	}

	return &Driver{client: client}
}

// Open creates a new Redis store with a specified instance.
func (s *Driver) Open(args ...interface{}) driver.Driver {
	return Open(args...)
}

// Count returns numbers of keys in store.
func (s *Driver) Count() int64 {
	return s.client.DbSize().Val()
}

// Exists returns true when a key exists false when not existing in store.
func (s *Driver) Exists(key string) bool {
	return s.client.Exists(key).Val()
}

// Get value from key in store.
func (s *Driver) Get(key string) ([]byte, error) {
	res, err := s.client.Get(key).Result()

	if len(res) == 0 {
		return nil, err
	}

	return []byte(res), err
}

// Set key with value in store.
func (s *Driver) Set(key string, value []byte) error {
	return s.client.Set(key, value, 0).Err()
}

// Delete key from store.
func (s *Driver) Delete(key string) error {
	return s.client.Del(key).Err()
}

// Close does not exists for Redis driver.
func (s *Driver) Close() error {
	return nil
}
