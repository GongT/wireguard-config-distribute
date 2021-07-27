package vpnManager

import (
	"fmt"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type vpnConfig struct {
	id types.VpnIdType

	Prefix               string                   `json:"prefix"`
	Allocations          map[string]NumberBasedIp `json:"allocations"`
	WireguardPrivateKeys map[string]string        `json:"wireguard-keys"`
	DefaultMtu           uint32                   `json:"mtu"`
	EnableObfuse         bool                     `json:"obfuse"`

	reAllocations   map[NumberBasedIp]bool
	prefixFreeParts uint
}

func (vpn *vpnConfig) getHostDomain() string {
	return vpn.id.Serialize()
}

func (vpn *vpnConfig) allocate(hostname string, requestIp NumberBasedIp) {
	tools.Error("allocate address %s.%s to client %s", vpn.Prefix, requestIp.String(vpn.prefixFreeParts), hostname)
	vpn.Allocations[hostname] = requestIp
	vpn.reAllocations[requestIp] = true
}

func (vpn *vpnConfig) format(hostname string) string {
	return vpn.Prefix + "." + vpn.Allocations[hostname].String(vpn.prefixFreeParts)
}

func (vpn *vpnConfig) cacheAndNormalize() {
	vpn.reAllocations = make(map[NumberBasedIp]bool)
	for _, ip := range vpn.Allocations {
		vpn.reAllocations[ip] = true
	}
	if vpn.DefaultMtu < MIN_VALID_MTU {
		vpn.DefaultMtu = DEFAULT_MTU
	}
}

func (vpn *vpnConfig) calcAllocSpace() error {
	fp := (3 - strings.Count(vpn.Prefix, "."))
	if fp < 1 {
		return fmt.Errorf("ip [%s] should have space to allocate", vpn.Prefix)
	}
	vpn.prefixFreeParts = uint(fp)
	return nil
}
