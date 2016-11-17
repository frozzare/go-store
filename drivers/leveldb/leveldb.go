package leveldb

import (
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

	if len(s.args) > 0 {
		path = s.args[0].(string)
	}

	if len(s.args) > 1 {
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
func (s *Driver) Get(key string) (value []byte, err error) {
	defer s.Close()

	return s.db().Get([]byte(key), nil)
}

// Set key with value in store.
func (s *Driver) Set(key string, value []byte) error {
	defer s.Close()

	return s.db().Put([]byte(key), value, nil)
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
