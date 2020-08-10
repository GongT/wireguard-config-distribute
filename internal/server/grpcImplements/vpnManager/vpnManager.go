package vpnManager

import (
	"errors"
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

			if err := vpn.calcAllocSpace(); err != nil {
				tools.Die("invalid config: VPN %s wrong prefix: %s", name, err.Error())
			}
			vpn.cache()
		}
	} else {
		add(mapper, storage, "default", &vpnConfig{
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

func add(mapper map[string]*vpnConfig, storage *storage.ServerStorage, name string, config *vpnConfig) error {
	if _, ok := mapper[name]; ok {
		return errors.New("Adding vpn name is already exists")
	}

	if err := config.calcAllocSpace(); err != nil {
		return err
	}
	config.cache()
	mapper[name] = config

	return nil
}

func (vpns *VpnManager) saveFile() error {
	return vpns.storage.WriteJson(VPN_STORE_NAME, vpns.mapper)
}

func (vpns *VpnManager) AddVpnSpace(name string, config vpnConfig) error {
	vpns.m.Lock()
	defer vpns.m.Unlock()

	return add(vpns.mapper, vpns.storage, name, &config)
}

func (vpns *VpnManager) GetLocked(name string) (*VpnHelper, bool) {
	vpn, ok := vpns.mapper[name]
	if ok {
		return createHelper(vpns, vpn, name), true
	} else {
		return nil, false
	}
}
