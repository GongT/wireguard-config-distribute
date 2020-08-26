package transport

import (
	"net"
	"strconv"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func getFree() (uint16, error) {
	lis, err := net.ListenUDP("udp", nil)
	if err != nil {
		return 0, err
	}
	port := lis.LocalAddr().(*net.UDPAddr).Port
	lis.Close()

	return uint16(port), nil
}

func format(ip string, port uint16) string {
	if strings.Contains(ip, ":") {
		return "[" + ip + "]:" + strconv.FormatUint(uint64(port), 10)
	} else {
		return ip + ":" + strconv.FormatUint(uint64(port), 10)
	}
}

func parse(ip string, port uint16) *net.UDPAddr {
	addrStr := format(ip, port)
	remoteAddr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		tools.Die("failed parse address: %v", addrStr)
	}
	return remoteAddr
}

func isSocketClosed(err error) bool {
	return strings.Contains(err.Error(), " closed ")
}
