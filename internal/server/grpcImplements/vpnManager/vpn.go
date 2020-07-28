package vpnManager

import (
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type vpnConfig struct {
	Prefix      string                   `json:"prefix"`
	Allocations map[string]NumberBasedIp `json:"allocations"`

	reAllocations   map[NumberBasedIp]bool
	prefixFreeParts uint
}

func (vpn *vpnConfig) allocate(hostname string, requestIp NumberBasedIp) {
	tools.Error("allocate address %s.%s to client %s", vpn.Prefix, requestIp.String(vpn.prefixFreeParts), hostname)
	vpn.Allocations[hostname] = requestIp
	vpn.reAllocations[requestIp] = true
}

func (vpn *vpnConfig) format(hostname string) string {
	return vpn.Prefix + "." + vpn.Allocations[hostname].String(vpn.prefixFreeParts)
}
