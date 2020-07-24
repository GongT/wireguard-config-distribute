package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func listenUnix(socketPath string) *net.Listener {
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		tools.Die("failed to listen! %s", err.Error())
	}
	fmt.Printf("Server listen on: %s\n", lis.Addr().String())
	return &lis
}

func listenTCP(port uint16) *net.Listener {
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.FormatInt(int64(port), 10))
	if err != nil {
		tools.Die("failed to listen! %s", err.Error())
	}
	fmt.Printf("Server listen on: %s\n", lis.Addr().String())
	return &lis
}
