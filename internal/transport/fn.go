package transport

import (
	"fmt"
	"net"
	"strconv"
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
	return ip + ":" + strconv.FormatUint(uint64(port), 10)
}

func parse(ip string, port uint16) *net.UDPAddr {
	addrStr := format(ip, port)
	remoteAddr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		panic(fmt.Errorf("failed parse address: %v", addrStr))
	}
	return remoteAddr
}
