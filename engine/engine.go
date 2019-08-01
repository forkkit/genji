package engine

import (
	"errors"
)

// Common errors returned by the engine implementations.
var (
	// ErrTransactionReadOnly must be returned when attempting to call write methods on a read-only transaction.
	ErrTransactionReadOnly = errors.New("transaction is read-only")

	// ErrStoreNotFound is returned when the targeted store doesn't exist.
	ErrStoreNotFound = errors.New("store not found")

	// ErrStoreAlreadyExists must be returned when attempting to create a store with the
	// same name as an existing one.
	ErrStoreAlreadyExists = errors.New("store already exists")

	// ErrKeyNotFound is returned when the targeted key doesn't exist.
	ErrKeyNotFound = errors.New("key not found")
)

// An Engine is responsible for storing data.
// Implementations can choose to store data on disk, in memory, in the browser etc. using the algorithms
// and data structures of their choice.
// Engines must support read-only and read/write transactions.
type Engine interface {
	Begin(writable bool) (Transaction, error)
	Close() error
}

// A Transaction provides methods for managing the collection of stores and the transaction itself.
// The transaction is either read-only or read/write. Read-only transactions can be used to read stores
// and read/write ones can be used to read, create, delete and modify stores.
type Transaction interface {
	Rollback() error
	Commit() error
	Store(name string) (Store, error)
	CreateStore(name string) error
	DropStore(name string) error
	StoreList(prefix string) ([]string, error)
}

// A Store manages key value pairs. It is an abstraction on top of data structures that can provide random readThe store can be implemented by any data stru
type Store interface {
	// Get returns a value associated with the given key. If no key is not found, it returns ErrKeyNotFound.
	Get(k []byte) ([]byte, error)
	// Put stores a key value pair. If it already exists, it overrides it.
	Put(k, v []byte) error
	// Delete a key value pair. If the key is not found, returns ErrKeyNotFound.
	Delete(k []byte) error
	// Truncate deletes all the key value pairs from the store.
	Truncate() error
	// AscendGreater seeks for the pivot and then goes through all the subsequent key value pairs in increasing order and calls the given function for each pair.
	// If the given function returns an error, the iteration stops and returns that error.
	// If the pivot is nil, starts from the beginning.
	AscendGreaterOrEqual(pivot []byte, fn func(k, v []byte) error) error
	// DescendGreater seeks for the pivot and then goes through all the subsequent key value pairs in descreasing order and calls the given function for each pair.
	// If the given function returns an error, the iteration stops and returns that error.
	// If the pivot is nil, starts from the end.
	DescendLessOrEqual(pivot []byte, fn func(k, v []byte) error) error
}
