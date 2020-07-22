package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc"
)

func (s *clientStateHolder) startNetwork() {
	for i := 0; i < 5; i++ {
		fmt.Printf("Connect to server: %s (try %d)\n", s.address, i)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, s.address, s.tlsOption, grpc.WithBlock())

		if err == nil {
			fmt.Println("  * grpc connect ok.")
			s.grpcConn = conn
			s.rpc = protocol.NewWireguardApiClient(conn)

			return
		} else {
			tools.Error("Failed connect server: %s", err.Error())
		}
	}
	tools.Die("Failed to connect server (after 5 retry).")
}
