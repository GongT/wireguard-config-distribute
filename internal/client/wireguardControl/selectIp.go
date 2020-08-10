package wireguardControl

import (
	"net"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	fastping "github.com/tatsushid/go-fastping"
)

const MAX_TRY = 3

var knownReachableIp map[string]bool = make(map[string]bool)
var knownUnreachableIp map[string]uint8 = make(map[string]uint8)

func selectIp(ips []string, v4Only bool) string {
	for _, ip := range ips {
		if knownReachableIp[ip] {
			return ip
		}
		if t, ok := knownUnreachableIp[ip]; ok && t >= MAX_TRY {
			return ""
		}
	}

	ipsFilter := make([]string, 0, len(ips))
	for _, ip := range ips {
		if t, ok := knownUnreachableIp[ip]; !ok || t < MAX_TRY {
			ipsFilter = append(ipsFilter, ip)
		}
	}

	if v4Only {
		ff := make([]string, 0, len(ipsFilter))
		for _, ip := range ipsFilter {
			if detect_ip.IsValidIPv4(ip) {
				ff = append(ff, ip)
			}
		}
		ipsFilter = ff
	}

	if len(ipsFilter) == 1 {
		return ipsFilter[0]
	}
	if len(ipsFilter) == 0 {
		return ""
	}
	ip := _selectIp(ipsFilter)
	if len(ip) > 0 {
		knownReachableIp[ip] = true
	}
	return ip
}

func _selectIp(ips []string) string {
	// link local
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		ipps := make([]net.IP, len(ips))
		for i, ip := range ips {
			ipps[i] = net.ParseIP(ip)
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				for _, ip := range ipps {
					if ipnet.Contains(ip) {
						return ip.String()
					}
				}
			}
		}
	} else {
		tools.Error("Failed list interfaces: %s", err.Error())
	}

	// ping each one
	ch := make(chan string, 2)

	p := fastping.NewPinger()

	for _, ip := range ips {
		p.AddIP(ip)
	}
	p.MaxRTT = 5 * time.Second
	p.OnRecv = func(addr *net.IPAddr, _ time.Duration) {
		ch <- addr.String()
	}
	p.OnIdle = func() {
		ch <- ""
	}

	go func() {
		tools.Debug("pinging %d address...", len(ips))
		err = p.Run()
		if err != nil {
			tools.Error("Failed run fastping: %s", err.Error())
		}
	}()

	select {
	case ret := <-ch:
		tools.Debug("pong: [%s]", ret)

		if len(ret) == 0 {
			for _, ip := range ips {
				if _, ok := knownUnreachableIp[ip]; ok {
					knownUnreachableIp[ip] += 1
				} else {
					knownUnreachableIp[ip] = 1
				}
			}
		}

		return ret
	}
}
