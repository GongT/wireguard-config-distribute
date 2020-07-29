package detect_ip

import (
	"net"
	"strings"
)

func IsValidIPv6(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}

	return strings.Contains(s, ":")
}

func IsValidIPv4(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}

	return !strings.Contains(s, ":")
}
