package boltdb

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"

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
func (s *Driver) db() (*bolt.DB, error) {
	if s.client != nil && !s.closed {
		return s.client, nil
	}

	s.closed = false

	path := "/tmp/store-bolt.db"
	mode := os.FileMode(0600)
	options := &bolt.Options{}

	if len(s.args) > 0 && s.args[0] != nil {
		path = s.args[0].(string)
	}

	if len(s.args) > 1 && s.args[1] != nil {
		mode = s.args[1].(os.FileMode)
	}

	if len(s.args) > 2 && s.args[2] != nil {
		options = s.args[2].(*bolt.Options)
	}

	client, err := bolt.Open(path, mode, options)

	if err != nil {
		return nil, err
	}

	s.client = client

	if len(s.bucket) == 0 {
		s.bucket = fmt.Sprintf("%x", md5.Sum([]byte(path)))
	}

	return client, nil
}

// Open creates a new BoltDB store.
func Open(args ...interface{}) (driver.Driver, error) {
	return &Driver{args: args}, nil
}

// Open creates a new BoltDB store with a specified instance.
func (s *Driver) Open(args ...interface{}) (driver.Driver, error) {
	return Open(args...)
}

// Count returns numbers of keys in store.
func (s *Driver) Count() (count int64, err error) {
	defer s.Close()

	db, err := s.db()

	if err != nil {
		return 0, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(s.bucket))
		}

		count = int64(bucket.Stats().KeyN)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return
}

// Exists returns true when a key exists false when not existing in store.
func (s *Driver) Exists(key string) (bool, error) {
	value, err := s.Get(key)

	if err != nil {
		return false, err
	}

	return value != nil, nil
}

// Keys returns a string slice with all keys.
func (s *Driver) Keys() (keys []string, err error) {
	defer s.Close()

	db, err := s.db()

	if err != nil {
		return []string{}, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(s.bucket))
		}

		bucket.ForEach(func(key []byte, value []byte) error {
			keys = append(keys, string(key))

			return nil
		})

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return
}

// Get returns the value for a key if any.
func (s *Driver) Get(key string, args ...interface{}) (value interface{}, err error) {
	defer s.Close()

	db, err := s.db()

	if err != nil {
		return nil, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(s.bucket))
		}

		var buffer bytes.Buffer
		buffer.Write(bucket.Get([]byte(key)))

		bytes := buffer.Bytes()

		var js interface{}

		if len(args) > 0 {
			js = args[0]
		}

		if err = json.Unmarshal(bytes, &js); err == nil {
			if len(args) > 0 {
				value = nil
			} else {
				value = js
			}
		} else if len(bytes) > 0 {
			value = string(bytes)
		}

		return nil
	})

	return
}

// Set key with value in store.
func (s *Driver) Set(key string, value interface{}) error {
	defer s.Close()

	db, err := s.db()

	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return err
		}

		if reflect.TypeOf(value).Kind() != reflect.String {
			value, err := json.Marshal(value)

			if err != nil {
				return err
			}

			err = bucket.Put([]byte(key), value)
		} else {
			err = bucket.Put([]byte(key), []byte(value.(string)))
		}

		if err != nil {
			return err
		}

		return nil
	})
}

// Delete key from store.
func (s *Driver) Delete(key string) error {
	defer s.Close()

	db, err := s.db()

	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
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
	db, err := s.db()

	if err != nil {
		return err
	}

	err = db.Close()

	if err != nil {
		return err
	}

	s.closed = true

	return nil
}

// Flush will remove all keys and values from the store.
func (s *Driver) Flush() error {
	defer s.Close()
	db, err := s.db()

	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(s.bucket))
	})
}
