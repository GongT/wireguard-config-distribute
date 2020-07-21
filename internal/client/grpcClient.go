package client

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (s *Client) ConnectGrpc() {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: config.IsDevelopmennt()})
	// remember to update address to use the new NGINX listen port
	address := config.GetConfig(config.CONFIG_SERVER_ADDRESS, config.CONFIG_SERVER_ADDRESS_DEFAULT)
	address = address + strconv.FormatInt(config.GetConfigNumber(config.CONFIG_SERVER_PORT, config.CONFIG_SERVER_PORT_DEFAULT), 10)

	fmt.Printf("Connecting to server: %s\n", address)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}

	s.grpcConn = conn
}
