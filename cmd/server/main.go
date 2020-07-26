package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	serverInternals "github.com/gongt/wireguard-config-distribute/internal/server"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements"
	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc/credentials"
)

var opts serverProgramOptions

func preparePassword(store *storage.ServerStorage) {
	save := func() {
		if err := store.WriteFile("password.txt", opts.Password); err != nil {
			tools.Die("Failed write password file: %s", err.Error())
		}
		tools.Error("Write password file")
	}

	savedPassword, _ := store.ReadFile("password.txt")
	savedPassword = strings.TrimSpace(savedPassword)

	noSet := len(opts.Password) == 0
	noExists := len(savedPassword) == 0

	if noSet {
		if noExists {
			opts.Password = tools.RandString(16)
			save()
		} else {
			opts.Password = savedPassword
			noSet = noExists
		}
	} else if opts.Password != savedPassword {
		save()
	}
}

func main() {
	config.InitProgramArguments(&opts)

	storagePath := opts.GetStorageLocation()
	if len(storagePath) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			tools.Die("Failed get user HOME: %s", err.Error())
		}
		storagePath = filepath.Join(home, ".wireguard-config-server")
		opts.StorageLocation = storagePath
	}
	store := storage.CreateStorage(storagePath)
	fmt.Printf("Storage path: %s\n", storagePath)

	preparePassword(store)

	var certs *credentials.TransportCredentials = nil
	if opts.GetGrpcInsecure() {
		if len(opts.GetGrpcServerKey()) > 0 || len(opts.GetGrpcServerPub()) > 0 {
			tools.Die("Can not use server-key/pub file with --insecure")
		}

		fmt.Println("Using insecure transport")
	} else {
		_certs, err := store.LoadOrCreateTLS(opts)
		if err != nil {
			tools.Die("Failed create or load TLS keyfile: %s", err.Error())
		}
		fmt.Println("Using TLS transport")
		certs = &_certs
	}

	server := serverInternals.NewServer(opts, certs, grpcImplements.CreateServerImplement(opts))

	server.Listen(opts)

	<-tools.WaitForCtrlC()

	server.Stop()
}
