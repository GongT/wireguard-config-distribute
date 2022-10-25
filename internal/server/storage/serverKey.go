package storage

import (
	"errors"
	"fmt"
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

var Ipv4LinkLocal = net.IPv4(224, 0, 0, 1)

func parseAll(ips []string) []net.IP {
	ret := make([]net.IP, 0, len(ips))
	for _, ipstr := range ips {
		if ip := net.ParseIP(ipstr); ip != nil {
			ret = append(ret, ip)
		}
	}
	return ret
}

func (storage *ServerStorage) loadOrCreateServerKey(options tlsOptions) error {
	if storage._cacheCaPri == nil || storage._cacheCa == nil {
		return errors.New("Invalid program state [CA Cache]")
	}

	ips := parseAll(options.GetPublicIp())

	var ipv4 net.IP
	detect_ip.RunDetect(&ipv4, &wrapGetIpOptions{options})
	if tools.IsIPv4(ipv4) {
		ips = append(ips, ipv4)
	}

	ips = append(ips, Ipv4LinkLocal, net.IPv6loopback)
	ips = append(ips, detect_ip.ListAllLocalNetworkIp()...)

	// unique.
	for j := len(ips) - 1; j >= 0; j-- {
		for i := j - 1; i >= 0; i-- {
			if ips[i].Equal(ips[j]) {
				ips = append(ips[:j], ips[j+1:]...)
				break
			}
		}
	}
	// unique end

	return storage.createServerKey(ips)
}

func (storage *ServerStorage) createServerKey(ipList []net.IP) error {
	var ipArr []net.IP
	fmt.Printf("Server addresses:\n")
	for _, ip := range ipList {
		fmt.Printf("  - %s\n", ip.String())
		ipArr = append(ipArr, ip)
	}

	return storage.createKey(ipArr, storage.PubFilePath(), storage.keyFilePath())
}

func (storage *ServerStorage) PubFilePath() string {
	return storage.Path("service.pem")
}

func (storage *ServerStorage) keyFilePath() string {
	return storage.Path("service.key")
}
