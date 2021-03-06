package driver

// Driver is the interface that must be implemented
// by a store driver.
type Driver interface {
	// Count counts the number of keys in the store.
	Count() (int64, error)

	// Delete delets a key and value from store if any.
	Delete(key string) error

	// Exists checks if a key exists in the store.
	Exists(key string) (bool, error)

	// Get returns the value for a key if any.
	Get(key string, args ...interface{}) (interface{}, error)

	// Keys returns a string slice with all keys.
	Keys() ([]string, error)

	// Open opens a new store.
	Open(args ...interface{}) (Driver, error)

	// Set key value in store.
	Set(key string, value interface{}) error

	// Close will close the driver connection if needed.
	Close() error

	// Flush will remove all keys and values from the store.
	Flush() error
}
