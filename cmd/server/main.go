package main

import (
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	serverInternals "github.com/gongt/wireguard-config-distribute/internal/server"
	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc"
)

func main() {
	opts := serverProgramOptions{}
	config.ParseProgramArguments(&opts)

	storagePath := opts.GetStorageLocation()
	if len(storagePath) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			tools.Die("Failed get user HOME: %s", err.Error())
		}
		storagePath = filepath.Join(home, ".wireguard-config-server")
	}
	store := storage.CreateStorage(storagePath)

	var transport grpc.ServerOption
	if opts.GetGrpcInsecure() {
		if len(opts.GetGrpcServerKey()) > 0 || len(opts.GetGrpcServerPub()) > 0 {
			tools.Die("Can not use server-key/pub file with --insecure")
		}

		transport = grpc.EmptyServerOption{}
	} else {
		certs, err := store.CreateTLSFilesIfNot(opts.GetGrpcServerKey(), opts.GetGrpcServerPub(), opts.GetServerName())
		if err != nil {
			tools.Die("Failed load TLS keyfile: %s", err.Error())
		}
		transport = grpc.Creds(certs)
	}

	server := serverInternals.NewServer(transport)

	server.Listen(opts)

	<-tools.WaitForCtrlC()

	server.Stop()
}
