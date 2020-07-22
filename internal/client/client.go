package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/config"
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

func NewClient() (client clientStateHolder) {
	client = clientStateHolder{}

	address := config.GetConfig(config.CONFIG_SERVER_ADDRESS, config.CONFIG_SERVER_ADDRESS_DEFAULT)
	address = address + ":" + strconv.FormatInt(config.GetConfigNumber(config.CONFIG_SERVER_PORT, config.CONFIG_SERVER_PORT_DEFAULT), 10)
	client.address = address

	// if config.IsDevelopmennt() {
	// 	fmt.Println("TLS did not enabled.")
	// 	client.tlsOption = grpc.WithInsecure()
	// } else {
	fmt.Println("TLS enabled.")
	creds := credentials.NewTLS(&tls.Config{})
	client.tlsOption = grpc.WithTransportCredentials(creds)
	// }

	client.quitChan = make(chan bool, 1)
	client.isQuit = false

	client.context = metadata.NewOutgoingContext(context.Background(), map[string][]string{})

	return client
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
