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

type vpnConfig struct {
	Prefix      string                   `json:"prefix"`
	Allocations map[string]NumberBasedIp `json:"allocations"`

	reAllocations   map[NumberBasedIp]bool
	prefixFreeParts uint
}

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

		for _, vpn := range mapper {
			if vpn.Allocations == nil {
				vpn.Allocations = make(map[string]NumberBasedIp)
			}
			if vpn.reAllocations == nil {
				vpn.reAllocations = make(map[NumberBasedIp]bool)
			}

			fp := (3 - strings.Count(vpn.Prefix, "."))
			if fp <= 1 {
				tools.Die("Invalid Config: VPN %s should have ip address space to allocate")
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
	err := storage.WriteJson(VPN_STORE_NAME, mapper)

	if err != nil {
		delete(mapper, name)
	}

	return err

}

func (vpns *VpnManager) Add(name string, config vpnConfig) error {
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

func (vpn *vpnConfig) allocate(hostname string, requestIp NumberBasedIp) {
	tools.Error("allocate address %s.%s to client %s", vpn.Prefix, requestIp.String(vpn.prefixFreeParts), hostname)
	vpn.Allocations[hostname] = requestIp
	vpn.reAllocations[requestIp] = true
}

func (vpn *vpnConfig) format(hostname string) string {
	return vpn.Prefix + "." + vpn.Allocations[hostname].String(vpn.prefixFreeParts)
}

func (vpns *VpnManager) AllocateIp(name string, hostname string, requestIp string) string {
	vpns.m.Lock()
	defer vpns.m.Unlock()

	vpn, ok := vpns.mapper[name]
	if !ok {
		tools.Die("VPN name %s must exists, but infact not.", name)
	}
	if vpn.reAllocations == nil {
		tools.Die("VPN staus %s.reAllocations must not nil.", name)
	}
	if vpn.Allocations == nil {
		tools.Die("VPN staus %s.Allocations must not nil.", name)
	}

	if _, exists := vpn.Allocations[hostname]; exists {
		return vpn.format(hostname)
	}

	reqIp := FromNumber(requestIp)
	if reqIp == 0 {
		reqIp = 1
	} else {
		if validRequest := net.ParseIP(vpn.Prefix + "." + requestIp); validRequest == nil {
			// request not valid
		} else if name, used := vpn.reAllocations[reqIp]; used {
			tools.Error("client %s want address %s, but used by %s", hostname, requestIp, name)
		} else {
			vpn.allocate(hostname, reqIp)
			return vpn.format(hostname)
		}
	}

	avaiable := NumberBasedIp(math.Pow(255.0, float64(reqIp)))
	for i := reqIp; i < avaiable; i += 1 {
		if _, used := vpn.reAllocations[reqIp]; !used {
			vpn.allocate(hostname, i)
			return vpn.format(hostname)
		}
	}

	return ""
}
