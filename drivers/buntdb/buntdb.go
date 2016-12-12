package buntdb

import (
	"encoding/json"
	"log"
	"reflect"

	bunt "github.com/tidwall/buntdb"

	"github.com/frozzare/go-store/driver"
)

// Driver represents a BundDB driver.
type Driver struct {
	args   []interface{}
	closed bool
	client *bunt.DB
}

// db returns the BundDB client if existing
// or creating a new if closed.
func (s *Driver) db() *bunt.DB {
	if s.client != nil && !s.closed {
		return s.client
	}

	s.closed = false

	path := "/tmp/store-bunt.db"

	if len(s.args) > 0 && s.args[0] != nil {
		path = s.args[0].(string)
	}

	db, err := bunt.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	s.client = db

	return s.client
}

// Open creates a new BundDB store.
func Open(args ...interface{}) driver.Driver {
	return &Driver{args: args}
}

// Open creates a new BundDB store with a specified instance.
func (s *Driver) Open(args ...interface{}) driver.Driver {
	return Open(args...)
}

// Count returns numbers of keys in store.
func (s *Driver) Count() (count int64) {
	defer s.Close()

	err := s.db().View(func(tx *bunt.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			count++

			return true
		})
	})

	if err != nil {
		return 0
	}

	return
}

// Exists returns true when a key exists false when not existing in store.
func (s *Driver) Exists(key string) (exists bool) {
	defer s.Close()

	err := s.db().View(func(tx *bunt.Tx) error {
		v, err := tx.Get(key)

		if err != nil {
			return err
		}

		exists = len(v) > 0

		return nil
	})

	if err != nil {
		return false
	}

	return
}

// Keys returns a string slice with all keys.
func (s *Driver) Keys() (keys []string, err error) {
	defer s.Close()

	err = s.db().View(func(tx *bunt.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			keys = append(keys, key)

			return true
		})
	})

	if err != nil {
		return []string{}, err
	}

	return keys, nil
}

// Get returns the value for a key if any.
func (s *Driver) Get(key string, args ...interface{}) (value interface{}, err error) {
	defer s.Close()

	err = s.db().View(func(tx *bunt.Tx) error {
		val, err := tx.Get(key)

		if err != nil {
			return err
		}

		if len(args) > 0 {
			value = args[0]
		}

		if err = json.Unmarshal([]byte(val), &value); err == nil {
			return nil
		}

		value = val

		return nil
	})

	return
}

// Set key with value in store.
func (s *Driver) Set(key string, value interface{}) error {
	defer s.Close()

	return s.db().Update(func(tx *bunt.Tx) error {
		if reflect.TypeOf(value).Kind() != reflect.String {
			value, err := json.Marshal(value)

			if err != nil {
				return err
			}

			tx.Set(key, string(value), nil)
		} else {

			tx.Set(key, value.(string), nil)
		}

		return nil
	})
}

// Delete key from store.
func (s *Driver) Delete(key string) error {
	defer s.Close()

	return s.db().Update(func(tx *bunt.Tx) error {
		_, err := tx.Delete(key)

		return err
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

// Flush will remove all keys and values from the store.
func (s *Driver) Flush() error {
	defer s.Close()

	return s.db().Update(func(tx *bunt.Tx) error {
		return tx.DeleteAll()
	})
}
