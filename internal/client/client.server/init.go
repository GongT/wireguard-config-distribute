package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServerStatus struct {
	tlsOption grpc.DialOption
	context   context.Context
	address   string

	rpc        protocol.WireguardApiClient
	connection *grpc.ClientConn
}

func NewGrpcClient(address string, tls TLSOptions) (ret ServerStatus) {
	if !strings.Contains(address, ":") {
		address += ":" + config.DEFAULT_PORT
	}
	ret.address = address

	creds, err := createClientTls(tls)
	if err != nil {
		tools.Die("Failed create TLS: %s", err.Error())
	}
	ret.tlsOption = grpc.WithTransportCredentials(creds)

	ret.context = metadata.NewOutgoingContext(context.Background(), map[string][]string{})

	return
}

func (stat *ServerStatus) Connect() {
	if stat.connection != nil {
		tools.Die("State error: connection already started")
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("Connect to server: %s (try %d)\n", stat.address, i)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, stat.address, stat.tlsOption, grpc.WithBlock())

		if err == nil {
			fmt.Println("  * grpc connect ok.")
			stat.connection = conn
			stat.rpc = protocol.NewWireguardApiClient(conn)

			return
		} else {
			tools.Error("Failed connect server: %s", err.Error())
		}
	}
	tools.Die("Failed to connect server (after 5 retry).")
}

func (stat *ServerStatus) Disconnect() {
	if err := stat.Close(); err != nil {
		tools.Error("Failed send close command: %s", err.Error())
	}
	if err := stat.connection.Close(); err != nil {
		tools.Error("Failed disconnect network: %s", err.Error())
	}
	fmt.Println("grpc closed.")
}
