package detect_ip

import "net"

func IsValidIPv6(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}

	return len(ip) == net.IPv6len
}
