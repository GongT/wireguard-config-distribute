package storage

import (
	"errors"
	"fmt"
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (storage *ServerStorage) loadOrCreateServerKey(options tlsOptions) error {
	if storage._cacheCaPri == nil || storage._cacheCa == nil {
		return errors.New("Invalid program state [CA Cache]")
	}

	ips := options.GetPublicIp()

	var ipv4, ipv6 string
	detect_ip.Detect(&ipv4, &ipv6, !options.GetIpHttpDsiable(), false)
	if detect_ip.IsValidIPv4(ipv4) {
		ips = append(ips, ipv4)
	}
	if detect_ip.IsValidIPv6(ipv6) {
		ips = append(ips, ipv6)
	}

	ips = append(ips, "127.0.0.1", "::1")
	ips = append(ips, detect_ip.DetectLocalNetwork()...)
	ips = tools.ArrayUnique(ips)

	return storage.createServerKey(ips)
}

func (storage *ServerStorage) createServerKey(ipList []string) error {
	var ipArr []net.IP
	fmt.Printf("Server addresses:\n")
	for _, ipstr := range ipList {
		if ip := net.ParseIP(ipstr); ip != nil {
			fmt.Printf("  - %s\n", ip.String())
			ipArr = append(ipArr, ip)
		} else {
			fmt.Printf("  x %s\n", ipstr)
		}
	}

	return storage.createKey(ipArr, storage.PubFilePath(), storage.keyFilePath())
}

func (storage *ServerStorage) PubFilePath() string {
	return storage.Path("service.pem")
}

func (storage *ServerStorage) keyFilePath() string {
	return storage.Path("service.key")
}
