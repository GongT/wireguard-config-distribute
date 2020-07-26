package storage

import (
	"errors"

	"google.golang.org/grpc/credentials"
)

type tlsOptions interface {
	GetServerName() string
	GetGrpcServerKey() string
	GetGrpcServerPub() string
	GetPublicIp() []string
	GetIpHttpDsiable() bool
}

func (storage *ServerStorage) LoadOrCreateTLS(options tlsOptions) (credentials.TransportCredentials, error) {
	pub := options.GetGrpcServerPub()
	pri := options.GetGrpcServerKey()

	if (len(pub) > 0) != (len(pri) > 0) {
		return nil, errors.New("TLS public/private key must use together")
	}

	if len(pub) == 0 {
		var err error
		_, _, err = storage.loadOrCreateCA(options.GetServerName())
		if err != nil {
			return nil, err
		}

		err = storage.loadOrCreateServerKey(options)
		if err != nil {
			return nil, err
		}

		pub = storage.PubFilePath()
		pri = storage.keyFilePath()
	}

	return credentials.NewServerTLSFromFile(pub, pri)
}
