package wireguardControl

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type InterfaceOptions interface {
	GetListenPort() uint16
	GetInterfaceName() string
	GetMTU() uint16
}

type peerData struct {
	comment      string
	publicKey    string
	presharedKey string
	ip           string
	port         uint16
	keepAlive    uint
	privateIp    string
}

type PeersCache struct {
	peers      []peerData
	ifopts     InterfaceOptions
	configFile string

	mu sync.Mutex
}

func NewPeersCache(options InterfaceOptions) *PeersCache {
	dir, err := ioutil.TempDir("", "wireguard")
	if err != nil {
		log.Fatal(err)
	}

	return &PeersCache{
		peers:      make([]peerData, 20),
		ifopts:     options,
		configFile: filepath.Join(dir, options.GetInterfaceName()+".conf"),
	}
}

func (ic *PeersCache) CreatConfigFile() error {
	return ioutil.WriteFile(ic.configFile, []byte(ic.CreatConfig()), os.FileMode(0644))
}

func (ic *PeersCache) CreatConfig() string {
	return ""
}

func (ic *PeersCache) UpdatePeers(list []*protocol.Peers_Peer) {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	tools.Error("Updating peers:")
	ic.peers = ic.peers[0:0]
	for _, peer := range list {
		selectedIp := selectIp(peer.GetPeer().GetAddress())
		if len(selectedIp) == 0 {
			tools.Error("  * DROP <%s>, failed ping any of %v", peer.GetTitle(), peer.GetPeer().GetAddress())
			continue
		}

		tools.Error("  * <%s> %s -> %s", peer.GetMachineId(), peer.GetHostname(), selectedIp)
		ic.peers = append(ic.peers, peerData{
			comment:      peer.GetTitle(),
			publicKey:    peer.GetPeer().GetPublicKey(),
			presharedKey: "",
			ip:           selectedIp,
			port:         uint16(peer.GetPeer().GetPort()),
			keepAlive:    uint(peer.GetPeer().GetKeepAlive()),
			privateIp:    peer.GetPeer().GetVpnIp(),
		})
	}

	ic.updateInterface()
}

func (ic *PeersCache) updateInterface() error {
	err := ic.CreatConfigFile()
	if err != nil {
		return err
	}
	return update(ic.ifopts.GetInterfaceName(), ic.configFile)
}

/*
用法: C:\Program Files\WireGuard\wireguard.exe [
    (无参数)：提升并安装管理服务
    /installmanagerservice
    /installtunnelservice CONFIG_PATH
    /uninstallmanagerservice
    /uninstalltunnelservice TUNNEL_NAME
    /managerservice
    /tunnelservice CONFIG_PATH
    /ui CMD_READ_HANDLE CMD_WRITE_HANDLE CMD_EVENT_HANDLE LOG_MAPPING_HANDLE
    /dumplog OUTPUT_PATH
    /update [LOG_FILE]
]
*/
