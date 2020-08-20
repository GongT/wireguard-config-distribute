package vpnManager

import (
	"math"
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type VpnHelper struct {
	manager *VpnManager
	config  *vpnConfig
	name    string
}

func createHelper(vpns *VpnManager, config *vpnConfig, name string) *VpnHelper {
	vpns.m.Lock()

	return &VpnHelper{
		manager: vpns,
		config:  config,
		name:    name,
	}
}

func (helper *VpnHelper) Release() {
	helper.manager.m.Unlock()
}

func (helper *VpnHelper) Subnet() uint {
	return (4 - helper.config.prefixFreeParts) * 8
}

func (helper *VpnHelper) GetMTU(ifmtu uint32) uint32 {
	if ifmtu >= MIN_VALID_MTU {
		return ifmtu
	} else {
		return helper.config.DefaultMtu
	}
}

func (helper *VpnHelper) AllocateIp(hostname string, requestIp string) string {
	manager := helper.manager
	config := helper.config

	if config.reAllocations == nil {
		tools.Die("VPN staus %s.reAllocations must not nil.", helper.name)
	}
	if config.Allocations == nil {
		tools.Die("VPN staus %s.Allocations must not nil.", helper.name)
	}

	if _, exists := config.Allocations[hostname]; exists {
		return config.format(hostname)
	}

	var reqIp NumberBasedIp
	if len(requestIp) == 0 {
		reqIp = 1
	} else {
		reqIp = FromNumber(requestIp)
		if validRequest := net.ParseIP(config.Prefix + "." + requestIp); validRequest == nil {
			// request not valid
			reqIp = 1
		} else if name, used := config.reAllocations[reqIp]; used {
			tools.Error("client %s want address %s, but used by %s", hostname, requestIp, name)
		} else {
			config.allocate(hostname, reqIp)
			manager.saveFile()
			return config.format(hostname)
		}
	}

	maximum := NumberBasedIp(math.Pow(255.0, float64(config.prefixFreeParts)))
	for i := reqIp; i < maximum; i++ {
		if _, used := config.reAllocations[i]; !used {
			config.allocate(hostname, i)
			manager.saveFile()
			return config.format(hostname)
		}
	}

	tools.Debug("Failed alloc ip for %s, request=%d[%s], maximum=%d, size=%d", hostname, reqIp, requestIp, maximum, len(config.reAllocations))

	return ""
}
