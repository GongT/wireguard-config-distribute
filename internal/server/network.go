package server

import (
	"net"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func listenUnix(socketPath string) net.Listener {
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		tools.Die("failed to listen [[%s]]! %s", socketPath, err.Error())
	}
	return lis
}

func listenTCP(port uint16) net.Listener {
	tryLis := "0.0.0.0:" + strconv.FormatInt(int64(port), 10)
	lis, err := net.Listen("tcp", tryLis)
	if err != nil {
		tools.Die("failed to listen [[%s]]! %s", tryLis, err.Error())
	}
	return lis
}
