package driver

// Driver is the interface that must be implemented
// by a store driver.
type Driver interface {
	// Count counts the number of keys in the store.
	Count() int64

	// Delete delets a key and value from store if any.
	Delete(key string) error

	// Exists checks if a key exists in the store.
	Exists(key string) bool

	// Get returns the value for a key if any.
	Get(key string) ([]byte, error)

	// Open opens a new store.
	Open(args ...interface{}) Driver

	// Set key value in store.
	Set(key string, value []byte) error

	// Close will close the driver connection if needed.
	Close() error
}
