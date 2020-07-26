package server

import (
	"net"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func listenUnix(socketPath string) net.Listener {
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		tools.Die("failed to listen [[%s]]! %s", lis.Addr().String(), err.Error())
	}
	return lis
}

func listenTCP(port uint16) net.Listener {
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.FormatInt(int64(port), 10))
	if err != nil {
		tools.Die("failed to listen [[%s]]! %s", lis.Addr().String(), err.Error())
	}
	return lis
}
