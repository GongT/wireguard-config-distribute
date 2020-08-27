package clientType

import (
	"net"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	fastping "github.com/tatsushid/go-fastping"
)

const MAX_TRY = 3

type IpFilter uint8

const (
	DontFilter = 0
	NoIpV4     = 1
	NoIpV6     = 2
)

var knownReachableIp map[string]bool = make(map[string]bool)
var knownUnreachableIp map[string]uint8 = make(map[string]uint8)

func selectIp(ips []string, filter IpFilter) string {
	for _, ip := range ips {
		if knownReachableIp[ip] {
			tools.Debug("  -> known reachable: " + ip)
			return ip
		}
		if t, ok := knownUnreachableIp[ip]; ok && t >= MAX_TRY {
			tools.Debug("  -> known UN reachable: " + ip)
			return ""
		}
	}

	ipsFilter := make([]string, 0, len(ips))
	for _, ip := range ips {
		if filter != NoIpV4 && tools.IsValidIPv4(ip) {
			ipsFilter = append(ipsFilter, ip)
		} else if filter != NoIpV6 && tools.IsValidIPv6(ip) {
			ipsFilter = append(ipsFilter, ip)
		}
	}
	tools.Debug("  : filtered: %v", ipsFilter)

	if len(ipsFilter) == 1 {
		tools.Debug("  -> only one, force use")
		return ipsFilter[0]
	}
	if len(ipsFilter) == 0 {
		tools.Debug("  -> no ip usable")
		return ""
	}
	ip := _selectIp(ipsFilter)
	if len(ip) > 0 {
		tools.Debug("  -> selected: " + ip)
		knownReachableIp[ip] = true
	} else {
		tools.Debug("  -> select fail")
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
						tools.Debug("  : on link address found")
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
		if len(ret) == 0 {
			tools.Debug("  : unreachable")
			for _, ip := range ips {
				if _, ok := knownUnreachableIp[ip]; ok {
					knownUnreachableIp[ip] += 1
				} else {
					knownUnreachableIp[ip] = 1
				}
			}
		} else {
			tools.Debug("  : pong!")
		}

		return ret
	}
}
