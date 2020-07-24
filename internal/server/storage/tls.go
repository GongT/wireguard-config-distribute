package storage

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/credentials"
)

func (storage *ServerStorage) LoadOrCreateTLS(pub string, pri string, serverName string) (credentials.TransportCredentials, error) {
	if (len(pub) > 0) != (len(pri) > 0) {
		return nil, errors.New("TLS public/private key must use together")
	}

	if len(pub) == 0 {
		var err error
		err = storage.createTLS(serverName)
		if err != nil {
			return nil, err
		}

		pub = storage.PubFilePath()
		pri = storage.KeyFilePath()
	}

	return credentials.NewServerTLSFromFile(pub, pri)
}

func (storage *ServerStorage) createTLS(serverName string) (err error) {
	fmt.Println("Creating TLS keys...")

	caFile := storage.Path("ca.cert.pem")
	caPriFile := storage.Path("ca.key.pem")

	data, err := readCA()
	if err != nil {
		return
	}

	createCA(serverName, caFile, caPriFile)

	pub := storage.PubFilePath()
	pri := storage.KeyFilePath()

}

func (storage *ServerStorage) PubFilePath() string {
	return storage.Path("service.pem")
}

func (storage *ServerStorage) KeyFilePath() string {
	return storage.Path("service.key")
}
