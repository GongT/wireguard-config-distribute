package grpcImplements

import "github.com/gongt/wireguard-config-distribute/internal/server/storage"

type ServerImplementOptions interface {
	GetPassword() string
	GetStorageLocation() string
	GetGrpcInsecure() bool
}

type serverImplement struct {
	password string
	storage  *storage.ServerStorage
	insecure bool
}

func CreateServerImplement(opts ServerImplementOptions) *serverImplement {
	store := storage.CreateStorage(opts.GetStorageLocation())
	return &serverImplement{
		password: opts.GetPassword(),
		storage:  store,
		insecure: opts.GetGrpcInsecure(),
	}
}
