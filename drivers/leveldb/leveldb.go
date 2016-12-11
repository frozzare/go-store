package leveldb

import (
	"encoding/json"
	"reflect"

	"github.com/frozzare/go-store/driver"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Driver represents a BoltDB driver.
type Driver struct {
	args   []interface{}
	closed bool
	client *leveldb.DB
}

// db returns the LevelDB client if existing
// or creating a new if closed.
func (s *Driver) db() *leveldb.DB {
	if s.client != nil && !s.closed {
		return s.client
	}

	s.closed = false

	path := "/tmp/store.leveldb"
	var options *opt.Options

	if len(s.args) > 0 && s.args[0] != nil {
		path = s.args[0].(string)
	}

	if len(s.args) > 1 && s.args[1] != nil {
		options = s.args[1].(*opt.Options)
	}

	client, err := leveldb.OpenFile(path, options)

	if err != nil {
		panic(err)
	}

	s.client = client

	return client
}

// Open creates a new LevelDB store.
func Open(args ...interface{}) driver.Driver {
	return &Driver{args: args}
}

// Open creates a new LevelDB store with a specified instance.
func (s *Driver) Open(args ...interface{}) driver.Driver {
	return Open(args...)
}

// Count returns numbers of keys in store.
func (s *Driver) Count() (count int64) {
	defer s.Close()

	iter := s.db().NewIterator(nil, nil)

	for iter.Next() {
		count++
	}

	iter.Release()

	if err := iter.Error(); err != nil {
		return 0
	}

	return
}

// Exists returns true when a key exists false when not existing in store.
func (s *Driver) Exists(key string) (ret bool) {
	ret, err := s.db().Has([]byte(key), nil)

	if err != nil {
		return false
	}

	return
}

// Get value from key in store.
func (s *Driver) Get(key string) (interface{}, error) {
	defer s.Close()

	res, err := s.db().Get([]byte(key), nil)

	if err != nil {
		return nil, err
	}

	var value interface{}
	if err = json.Unmarshal([]byte(res), &value); err == nil {
		return value, nil
	}

	return string(res), nil
}

// Keys returns a string slice with all keys.
func (s *Driver) Keys() ([]string, error) {
	defer s.Close()

	iter := s.db().NewIterator(nil, nil)

	var keys []string

	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}

	iter.Release()

	if err := iter.Error(); err != nil {
		return []string{}, err
	}

	return keys, nil
}

// Set key with value in store.
func (s *Driver) Set(key string, value interface{}) error {
	defer s.Close()

	if reflect.TypeOf(value).Kind() != reflect.String {
		value, err := json.Marshal(value)

		if err != nil {
			return err
		}

		return s.db().Put([]byte(key), value, nil)
	}

	return s.db().Put([]byte(key), []byte(value.(string)), nil)
}

// Delete key from store.
func (s *Driver) Delete(key string) error {
	defer s.Close()

	return s.db().Delete([]byte(key), nil)
}

// Close will close the boltdb client.
func (s *Driver) Close() error {
	err := s.db().Close()

	if err != nil {
		return err
	}

	s.closed = true

	return nil
}

// Flush will remove all keys and values from the store.
func (s *Driver) Flush() error {
	defer s.Close()

	db := s.db()
	iter := db.NewIterator(nil, nil)

	for iter.Next() {
		db.Delete(iter.Key(), nil)
	}

	iter.Release()

	return nil
}
