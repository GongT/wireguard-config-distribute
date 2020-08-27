package tools

import (
	"net"
)

func IsValidIPv6(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}

	return IsIPv6(ip)
}

func IsValidIPv4(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}

	return IsIPv4(ip)
}

func IsIPv6(ip net.IP) bool {
	if ip == nil {
		return false
	}
	p4 := ip.To4()
	return len(p4) != net.IPv4len
}

func IsIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}
	p4 := ip.To4()
	return len(p4) == net.IPv4len
}
