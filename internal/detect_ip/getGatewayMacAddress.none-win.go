// +build !windows

package detect_ip

import (
	"net"

	"github.com/j-keck/arping"
)

func findGatewayAddr(gwIp net.IP) (string, error) {
	ret, _, err := arping.Ping(gwIp)
	return ret.String(), err
}
