package storage

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type ServerStorage struct {
	path string

	_cacheCa    *x509.Certificate
	_cacheCaPri *rsa.PrivateKey
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

func (storage *ServerStorage) Path(name string) string {
	return filepath.Join(storage.path, name)
}

func (storage *ServerStorage) WriteFile(file string, content string) error {
	f := storage.Path(file)
	if err := os.MkdirAll(filepath.Dir(f), os.FileMode(0755)); err != nil {
		return err
	}
	return ioutil.WriteFile(f, []byte(content), os.FileMode(0644))
}

func (storage *ServerStorage) WriteJson(file string, data interface{}) error {
	f := storage.Path(file)
	if err := os.MkdirAll(filepath.Dir(f), os.FileMode(0755)); err != nil {
		return err
	}
	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, bs, os.FileMode(0644))
}

func (storage *ServerStorage) ReadFile(file string) (string, error) {
	f := storage.Path(file)
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func (storage *ServerStorage) ReadJson(file string, result interface{}) error {
	f := storage.Path(file)
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, result)
}

func (storage *ServerStorage) PathExists(file string) bool {
	stat, err := os.Stat(storage.Path(file))
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}
