package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServerStatus struct {
	tlsOption grpc.DialOption
	address   string

	context       context.Context
	contextMeta   map[string][]string
	contextCancel context.CancelFunc

	rpc        protocol.WireguardApiClient
	connection *grpc.ClientConn
}

func NewGrpcClient(address string, tls TLSOptions) (ret ServerStatus) {
	ret.address = address

	creds, err := createClientTls(tls)
	if err != nil {
		tools.Die("Failed create TLS: %s", err.Error())
	}
	ret.tlsOption = grpc.WithTransportCredentials(creds)

	ret.context, ret.contextCancel = context.WithCancel(metadata.NewOutgoingContext(context.Background(), ret.contextMeta))

	return
}

func (stat *ServerStatus) Connect() {
	if stat.connection != nil {
		tools.Die("State error: rpc connection already started")
	}
	fmt.Printf("Connect to server: %s\n", stat.address)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, stat.address, stat.tlsOption, grpc.WithBlock(), grpc.WithReturnConnectionError())

	if err != nil {
		tools.Die("Failed to connect server: %s.", err.Error())
	}

	fmt.Println("  * grpc connect ok.")
	stat.connection = conn
	stat.rpc = protocol.NewWireguardApiClient(conn)

	return
}

func (stat *ServerStatus) Disconnect(sessionId uint64) {
	stat.contextCancel()

	if sessionId != 0 {
		tools.Error("Sending close command.")
		if err := stat.Close(sessionId); err != nil {
			tools.Error("Failed send close command: %s", err.Error())
		}
	}
	tools.Error("Disconnect network.")
	if err := stat.connection.Close(); err != nil {
		tools.Error("Failed disconnect network: %s", err.Error())
	}
	fmt.Println("grpc gracefull closed.")
}
