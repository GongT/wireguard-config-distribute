package storage

import (
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type ServerStorage struct {
	path string
}

func CreateStorage(location string) *ServerStorage {
	store := ServerStorage{
		path: location,
	}

	err := os.MkdirAll(location, os.FileMode(0755))
	if err != nil {
		tools.Die("Failed create storage: %s", err.Error())
	}

	return &store
}

func (storage ServerStorage) Path(name string) string {
	return filepath.Join(storage.path, name)
}
