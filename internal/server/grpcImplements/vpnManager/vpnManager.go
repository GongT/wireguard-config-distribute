package vpnManager

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

const VPN_STORE_NAME = "vpns.json"

type VpnManager struct {
	storage *storage.ServerStorage
	mapper  map[types.VpnIdType] /* vpn name, eg: default */ *vpnConfig

	m sync.Mutex
}

func NewVpnManager(storage *storage.ServerStorage) *VpnManager {
	mapper := make(map[types.VpnIdType]*vpnConfig, 0)

	if storage.PathExists(VPN_STORE_NAME) {
		if storage.ReadJson(VPN_STORE_NAME, &mapper) != nil {
			tools.Die("Invalid content: " + storage.Path(VPN_STORE_NAME))
		}

		for name, vpn := range mapper {
			if vpn.Allocations == nil {
				vpn.Allocations = make(map[string]NumberBasedIp)
			}

			if err := vpn.calcAllocSpace(); err != nil {
				tools.Die("invalid config: VPN %s wrong prefix: %s", name, err.Error())
			}
			vpn.cacheAndNormalize()
		}
	} else {
		add(mapper, storage, types.DeSerializeVpnIdType("default"), &vpnConfig{
			Prefix:      "10.166",
			Allocations: make(map[string]NumberBasedIp),
		})
	}

	ret := VpnManager{
		storage: storage,
		mapper:  mapper,
	}

	return &ret
}

func add(mapper map[types.VpnIdType]*vpnConfig, storage *storage.ServerStorage, name types.VpnIdType, config *vpnConfig) error {
	if _, ok := mapper[name]; ok {
		return errors.New("Adding vpn name is already exists")
	}
	if err := config.calcAllocSpace(); err != nil {
		return err
	}
	config.id = name
	config.cacheAndNormalize()
	mapper[name] = config

	return nil
}

func (vpns *VpnManager) saveFile() error {
	return vpns.storage.WriteJson(VPN_STORE_NAME, vpns.mapper)
}

func (vpns *VpnManager) AddVpnSpace(name types.VpnIdType, config vpnConfig) error {
	vpns.m.Lock()
	defer vpns.m.Unlock()

	return add(vpns.mapper, vpns.storage, name, &config)
}

func (vpns *VpnManager) GetLocked(name types.VpnIdType) (*VpnHelper, bool) {
	vpn, ok := vpns.mapper[name]
	if ok {
		return createHelper(vpns, vpn, name), true
	} else {
		return nil, false
	}
}

func (vpns *VpnManager) Dump() string {
	vpns.m.Lock()
	defer vpns.m.Unlock()

	ret := "== Vpn Mamager Status ==\n"

	ret += "Storage path: " + vpns.storage.Path(VPN_STORE_NAME) + "\n"

	for name, vpn := range vpns.mapper {
		ret += fmt.Sprintf("> %v: %v(%v) | MTU=%v | OBFS=%v\n", name, vpn.Prefix, vpn.prefixFreeParts, vpn.DefaultMtu, vpn.EnableObfuse)
		for host, ip := range vpn.Allocations {
			ret += fmt.Sprintf("\t%-15v => %v\n", vpn.Prefix+"."+ip.String(vpn.prefixFreeParts), host)
		}
	}

	return ret
}
