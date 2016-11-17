package boltdb

import (
	"bytes"
	"crypto/md5"
	"fmt"

	"os"

	"github.com/boltdb/bolt"
	"github.com/frozzare/go-store/driver"
)

// Driver represents a BoltDB driver.
type Driver struct {
	args   []interface{}
	bucket string
	closed bool
	client *bolt.DB
}

// db returns the BoltDB client if existing
// or creating a new if closed.
func (s *Driver) db() *bolt.DB {
	if s.client != nil && !s.closed {
		return s.client
	}

	s.closed = false

	path := "/tmp/store.db"
	mode := os.FileMode(0600)
	options := &bolt.Options{}

	if len(s.args) > 0 {
		path = s.args[0].(string)
	}

	if len(s.args) > 1 {
		mode = s.args[1].(os.FileMode)
	}

	if len(s.args) > 2 {
		options = s.args[2].(*bolt.Options)
	}

	client, err := bolt.Open(path, mode, options)

	if err != nil {
		panic(err)
	}

	s.client = client

	if len(s.bucket) == 0 {
		s.bucket = fmt.Sprintf("%x", md5.Sum([]byte(path)))
	}

	return client
}

// Open creates a new BoltDB store.
func Open(args ...interface{}) driver.Driver {
	return &Driver{args: args}
}

// Open creates a new BoltDB store with a specified instance.
func (s *Driver) Open(args ...interface{}) driver.Driver {
	return Open(args...)
}

// Count returns numbers of keys in store.
func (s *Driver) Count() (count int64) {
	defer s.Close()

	err := s.db().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(s.bucket))
		}

		count = int64(bucket.Stats().KeyN)
		return nil
	})

	if err != nil {
		return 0
	}

	return
}

// Exists returns true when a key exists false when not existing in store.
func (s *Driver) Exists(key string) bool {
	value, err := s.Get(key)

	if err != nil {
		return false
	}

	return len(value) > 0
}

// Get value from key in store.
func (s *Driver) Get(key string) (value []byte, err error) {
	defer s.Close()

	err = s.db().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(s.bucket))
		}

		var buffer bytes.Buffer
		buffer.Write(bucket.Get([]byte(key)))

		value = buffer.Bytes()
		return nil
	})

	return
}

// Set key with value in store.
func (s *Driver) Set(key string, value []byte) error {
	defer s.Close()

	return s.db().Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), value)
		if err != nil {
			return err
		}

		return nil
	})
}

// Delete key from store.
func (s *Driver) Delete(key string) error {
	defer s.Close()

	return s.db().Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return err
		}

		err = bucket.Delete([]byte(key))
		if err != nil {
			return err
		}

		return nil
	})
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