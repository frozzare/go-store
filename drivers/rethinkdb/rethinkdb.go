package rethinkdb

import (
	"encoding/json"
	"math/rand"

	"github.com/frozzare/go-store/driver"

	r "gopkg.in/gorethink/gorethink.v3"
)

// Driver represents a Redis driver.
type Driver struct {
	session *r.Session
	table   string
}

// Open creates a new Redis store.
func Open(args ...interface{}) (driver.Driver, error) {
	var options r.ConnectOpts
	var table string

	if len(args) > 0 && args[0] != nil {
		options = args[0].(r.ConnectOpts)
	} else {
		options = r.ConnectOpts{
			Address:  "localhost:28015",
			Database: "store",
		}
	}

	session, err := r.Connect(options)

	if err != nil {
		return nil, err
	}

	if len(args) > 1 {
		table = args[1].(string)
	} else {
		table = "store"
	}

	res, _ := r.DBCreate(options.Database).Run(session)

	defer res.Close()

	res, err = r.TableList().Run(session)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	var existing []string
	res.All(&existing)

	var found bool
	for _, item := range existing {
		if item == table {
			found = true
			break
		}
	}

	if !found {
		_, err = r.TableCreate(table).Run(session)
	}

	if err != nil {
		return nil, err
	}

	return &Driver{session: session, table: table}, nil
}

// Open creates a new Redis store with a specified instance.
func (s *Driver) Open(args ...interface{}) (driver.Driver, error) {
	return Open(args...)
}

// Count returns numbers of keys in store.
func (s *Driver) Count() (int64, error) {
	res, err := r.Table(s.table).Count().Run(s.session)

	defer res.Close()

	if err != nil {
		return 0, err
	}

	var rows []int

	res.All(&rows)

	return int64(rows[rand.Intn(len(rows))]), nil
}

// Exists returns true when a key exists false when not existing in store.
func (s *Driver) Exists(key string) (bool, error) {
	res, err := s.Get(key)

	return res != nil && err == nil, err
}

// Get returns the value for a key if any.
func (s *Driver) Get(key string, args ...interface{}) (interface{}, error) {
	res, err := r.Table(s.table).Get(key).Run(s.session)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	var row interface{}
	err = res.One(&row)

	if err != nil {
		return nil, err
	}

	value := row.(map[string]interface{})["value"]

	if len(args) == 0 {
		return value, nil
	}

	j, err := json.Marshal(value)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(j, &args[0]); err != nil {
		return nil, err
	}

	return nil, nil
}

// Keys returns a string slice with all keys.
func (s *Driver) Keys() ([]string, error) {
	res, err := r.Table(s.table).Run(s.session)

	defer res.Close()

	if err != nil {
		return []string{}, nil
	}

	var rows []map[string]interface{}

	res.All(&rows)

	var keys []string

	for _, row := range rows {
		keys = append(keys, row["id"].(string))
	}

	return keys, nil
}

// Set key with value in store.
func (s *Driver) Set(key string, value interface{}) error {
	_, err := r.Table(s.table).Insert(map[string]interface{}{
		"id":    key,
		"value": value,
	}, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(s.session)

	if err != nil {
		return err
	}

	return nil
}

// Delete key from store.
func (s *Driver) Delete(key string) error {
	_, err := r.Table(s.table).Get(key).Delete().RunWrite(s.session)

	if err != nil {
		return err
	}

	return nil
}

// Close will close the RethinkDB session.
func (s *Driver) Close() error {
	s.session.Close()

	return nil
}

// Flush will remove all keys and values from the store.
func (s *Driver) Flush() error {
	_, err := r.Table(s.table).Delete().RunWrite(s.session)

	if err != nil {
		return err
	}

	return nil
}
