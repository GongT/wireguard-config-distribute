package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type clientStateHolder struct {
	id        string
	rpc       protocol.WireguardApiClient
	grpcConn  *grpc.ClientConn
	address   string
	tlsOption grpc.DialOption
	quitChan  chan bool
	isQuit    bool
	context   context.Context
}

type whereToConnect interface {
	GetServer() string
}

func NewClient(options whereToConnect, creds credentials.TransportCredentials) clientStateHolder {
	c := clientStateHolder{}

	c.address = options.GetServer()

	c.tlsOption = grpc.WithTransportCredentials(creds)

	c.quitChan = make(chan bool, 1)
	c.isQuit = false

	c.context = metadata.NewOutgoingContext(context.Background(), map[string][]string{})

	return c
}
func (s *clientStateHolder) Quit() {
	if s.isQuit {
		tools.Error("Duplicate call to Client.quit()")
		return
	}
	s.isQuit = true

	err := s.grpcConn.Close()
	if err != nil {
		tools.Error("Failed disconnect grpc: %s", err.Error())
	}
	fmt.Println("grpc closed.")

	s.quitChan <- true
}

func (s *clientStateHolder) StartNetwork() {
	s.startNetwork()
	s.startProtocol()
}

func (s *clientStateHolder) StartCommunication() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		fmt.Println("start communication...")
		for {
			select {
			case <-ticker.C:
				_, err := s.rpc.KeepAlive(s.context, tools.EmptyPb)
				if err != nil {
					tools.Error("grpc keep alive failed, is server running? %s", err.Error())
				}
			case <-s.quitChan:
				ticker.Stop()
				fmt.Println("stop communication.")
				return
			}
		}
	}()
}
