package vpnManager

import (
	"errors"
	"math"
	"net"
	"strings"
	"sync"

	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

const VPN_STORE_NAME = "vpns.json"

type VpnManager struct {
	storage *storage.ServerStorage
	mapper  map[string]*vpnConfig

	m sync.Mutex
}

func NewVpnManager(storage *storage.ServerStorage) *VpnManager {
	mapper := make(map[string]*vpnConfig, 0)

	if storage.PathExists(VPN_STORE_NAME) {
		if storage.ReadJson(VPN_STORE_NAME, &mapper) != nil {
			tools.Die("Invalid content: " + storage.Path(VPN_STORE_NAME))
		}

		for name, vpn := range mapper {
			if vpn.Allocations == nil {
				vpn.Allocations = make(map[string]NumberBasedIp)
			}
			if vpn.reAllocations == nil {
				vpn.reAllocations = make(map[NumberBasedIp]bool)
			}

			fp := (3 - strings.Count(vpn.Prefix, "."))
			if fp < 1 {
				tools.Die("Invalid Config: VPN %s should have ip address space to allocate", name)
			}
			vpn.prefixFreeParts = uint(fp)

			for _, ip := range vpn.Allocations {
				vpn.reAllocations[ip] = true
			}
		}
	} else {
		add(mapper, storage, "default", &vpnConfig{
			Prefix:          "10.166",
			prefixFreeParts: 2,
			reAllocations:   make(map[NumberBasedIp]bool),
			Allocations:     make(map[string]NumberBasedIp),
		})
	}

	ret := VpnManager{
		storage: storage,
		mapper:  mapper,
	}

	return &ret
}

func add(mapper map[string]*vpnConfig, storage *storage.ServerStorage, name string, config *vpnConfig) error {
	if _, ok := mapper[name]; ok {
		return errors.New("Adding vpn name is already exists")
	}

	mapper[name] = config

	return nil
}

func (vpns *VpnManager) _save() error {
	return vpns.storage.WriteJson(VPN_STORE_NAME, vpns.mapper)
}

func (vpns *VpnManager) AddVpnSpace(name string, config vpnConfig) error {
	vpns.m.Lock()
	defer vpns.m.Unlock()

	return add(vpns.mapper, vpns.storage, name, &config)
}

func (vpns *VpnManager) Exists(name string) bool {
	vpns.m.Lock()
	defer vpns.m.Unlock()

	_, ok := vpns.mapper[name]
	return ok
}

func (vpns *VpnManager) AllocateIp(vpnname string, hostname string, requestIp string) string {
	vpns.m.Lock()
	defer vpns.m.Unlock()

	vpn, ok := vpns.mapper[vpnname]
	if !ok {
		tools.Die("VPN name %s must exists, but infact not.", vpnname)
	}
	if vpn.reAllocations == nil {
		tools.Die("VPN staus %s.reAllocations must not nil.", vpnname)
	}
	if vpn.Allocations == nil {
		tools.Die("VPN staus %s.Allocations must not nil.", vpnname)
	}

	if _, exists := vpn.Allocations[hostname]; exists {
		return vpn.format(hostname)
	}

	var reqIp NumberBasedIp
	if len(requestIp) == 0 {
		reqIp = 1
	} else {
		reqIp = FromNumber(requestIp)
		if validRequest := net.ParseIP(vpn.Prefix + "." + requestIp); validRequest == nil {
			// request not valid
			reqIp = 1
		} else if name, used := vpn.reAllocations[reqIp]; used {
			tools.Error("client %s want address %s, but used by %s", hostname, requestIp, name)
		} else {
			vpn.allocate(hostname, reqIp)
			vpns._save()
			return vpn.format(hostname)
		}
	}

	maximum := NumberBasedIp(math.Pow(255.0, float64(vpn.prefixFreeParts)))
	for i := reqIp; i < maximum; i++ {
		if _, used := vpn.reAllocations[i]; !used {
			vpn.allocate(hostname, i)
			vpns._save()
			return vpn.format(hostname)
		}
	}

	tools.Debug("Failed alloc ip for %s, request=%d[%s], maximum=%d, size=%d", hostname, reqIp, requestIp, maximum, len(vpn.reAllocations))

	return ""
}
