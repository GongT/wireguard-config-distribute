package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	serverInternals "github.com/gongt/wireguard-config-distribute/internal/server"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements"
	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
	"github.com/gongt/wireguard-config-distribute/internal/systemd"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc/credentials"
)

var opts *serverProgramOptions = &serverProgramOptions{}

func main() {
	spew.Config.Indent = "    "
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.MaxDepth = 3

	log.Println("program start.")
	if err := config.InitProgramArguments(opts); err != nil {
		tools.Die("invalid commandline arguments: %s", err.Error())
	}

	fmt.Printf("Storage path: %s\n", opts.GetStorageLocation())
	store := storage.CreateStorage(opts.GetStorageLocation())

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

	impl := grpcImplements.CreateServerImplement(opts)
	server := serverInternals.NewServer(opts, certs, impl)

	server.Listen(opts)
	impl.StartWorker()

	systemd.ChangeToReady()
	<-tools.WaitForCtrlC()
	systemd.ChangeToQuit()

	impl.Quit()
	server.Stop()
}

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
