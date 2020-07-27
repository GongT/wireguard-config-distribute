package networkManager

import (
	"sync"

	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

const NETWORK_STORE_NAME = "networks.json"

type NetworkManager struct {
	storage *storage.ServerStorage
	list    []string

	m sync.Mutex
}

func NewNetworkManager(storage *storage.ServerStorage) *NetworkManager {
	list := make([]string, 0)
	if storage.PathExists(NETWORK_STORE_NAME) {
		if storage.ReadJson(NETWORK_STORE_NAME, list) != nil {
			tools.Die("Invalid content: " + storage.Path(NETWORK_STORE_NAME))
		}
	} else {
		// networks.list=append(networks.list, defaultNetwork)
	}

	networks := NetworkManager{
		storage: storage,
		list:    list,
	}

	return &networks
}

func (networks *NetworkManager) Add(name string) error {
	networks.m.Lock()
	defer networks.m.Unlock()

	networks.list = append(networks.list, name)
	err := networks.storage.WriteJson(NETWORK_STORE_NAME, networks.list)

	if err != nil {
		networks.list = networks.list[:len(networks.list)-1]
	}

	return err
}

func (networks *NetworkManager) Exists(name string) bool {
	if name == "@alone" {
		return true
	}

	networks.m.Lock()
	defer networks.m.Unlock()

	return tools.ArrayContains(networks.list, name)
}
