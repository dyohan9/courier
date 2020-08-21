package courier

import (
	"context"
	"fmt"
	"strings"
)

type StorageConstructorFunc func(*Config) (Storage, error)

type Storage interface {
	// Test tests whether the storage is properly configured
	Test(ctx context.Context) error

	// PutFile writes the passed in file to the storage with the passed in content type
	PutFile(ctx context.Context, path string, contentType string, content []byte) (string, error)
}

// NewStorage creates the type of storage passed in
func NewStorage(config *Config) (Storage, error) {
	storageFunc, found := registeredStorages[strings.ToLower(config.Storage)]
	if !found {
		return nil, fmt.Errorf("no such storage type: '%s'", config.Storage)
	}
	return storageFunc(config)
}

// RegisterStorage adds a new storage, called by individual storages in their init() func
func RegisterStorage(storageType string, constructorFunc StorageConstructorFunc) {
	registeredStorages[strings.ToLower(storageType)] = constructorFunc
}

var registeredStorages = make(map[string]StorageConstructorFunc)
