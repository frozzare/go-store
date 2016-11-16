package store

import (
	"fmt"
	"sync"

	"github.com/frozzare/go-store/driver"
	"github.com/frozzare/go-store/drivers/redis"
	"github.com/frozzare/go-store/drivers/rwmutex"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

// init register the default store driver.
func init() {
	Register("rwmutex", &rwmutex.Driver{})
	Register("redis", &redis.Driver{})
}

// Register makes a store driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver driver.Driver) {
	driversMu.Lock()

	defer driversMu.Unlock()

	if driver == nil {
		panic("store: Register driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("store: Register called twice for driver " + name)
	}

	drivers[name] = driver
}

// Open opens a store driver and return it's implementation
// or a error if driver is not found.
func Open(args ...interface{}) (driver.Driver, error) {
	name := "rwmutex"
	if len(args) > 0 {
		name = args[0].(string)
		args = args[1:]
	}

	driversMu.RLock()
	driver, ok := drivers[name]
	driversMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("store: unknown driver %q (forgotten import?)", name)
	}

	return driver.Open(args...), nil
}
