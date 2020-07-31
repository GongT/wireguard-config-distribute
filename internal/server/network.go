package server

import (
	"net"
	"os"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func listenUnix(socketPath string) net.Listener {
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		tools.Die("failed to listen [[%s]]! %s", socketPath, err.Error())
	}
	err = os.Chmod(socketPath, os.FileMode(0777))
	if err != nil {
		tools.Die("failed to chmod socket file! %s", err.Error())
	}

	fi, err := os.Stat(socketPath)
	if err != nil {
		tools.Die("listen socket seems not work [[%s]]! %s", socketPath, err.Error())
	}
	tools.Debug("socket file mode:\n%s\t%s", socketPath, fi.Mode().String())

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
