package peerStatus

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
)

type PeerData struct {
	Title     string
	Hostname  string
	PublicKey string
	Address   string
	Port      uint32
	VpnIp     string
	KeepAlive uint32
	MTU       uint32
	Hosts     []string

	NetworkId    string
	ExternalIp   []string
	ExternalPort uint32
	InternalIp   []string
	InternalPort uint32
}

type PeerStatus struct {
	OnChange chan uint64
	list     map[uint64]*PeerData

	result protocol.Peers

	m sync.Mutex
}

func NewPeerStatus() *PeerStatus {
	return &PeerStatus{
		OnChange: make(chan uint64),
		list:     make(map[uint64]*PeerData),
	}
}

func (peers *PeerStatus) Delete(cid uint64) {
	peers.m.Lock()
	defer peers.m.Unlock()

	_, exists := peers.list[cid]

	if !exists {
		return
	}

	delete(peers.list, cid)
	peers.recreate()
	peers.OnChange <- cid
}

func (peers *PeerStatus) Add(cid uint64, peer *PeerData) {
	peers.m.Lock()
	defer peers.m.Unlock()

	old, exists := peers.list[cid]

	if exists && exactSame(old, peer) {
		return
	}

	peers.list[cid] = peer
	peers.recreate()
	peers.OnChange <- cid
}

func (peers *PeerStatus) recreate() {
	hosts := make(map[string]string)
	list := make([]*protocol.Peers_Peer, 0, len(peers.list))

	for _, peer := range peers.list {
		p := &protocol.Peers_Peer{
			Title:    peer.Title,
			Hostname: peer.Hostname,
			Peer: &protocol.Peers_ConnectionTarget{
				PublicKey: peer.PublicKey,
				Address:   peer.Address,
				Port:      peer.Port,
				VpnIp:     peer.VpnIp,
				KeepAlive: peer.KeepAlive,
				MTU:       peer.MTU,
			},
		}

		list = append(list, p)
		for _, host := range peer.Hosts {
			hosts[host] = peer.VpnIp
		}
	}

	peers.result = protocol.Peers{
		List:  list,
		Hosts: hosts,
	}
}

func (peers *PeerStatus) Serialize() *protocol.Peers {
	return &peers.result
}

func exactSame(a *PeerData, b *PeerData) bool {
	j1, err1 := json.Marshal(a)
	j2, err2 := json.Marshal(b)

	if err1 != nil || err2 != nil {
		return false
	}

	return bytes.Compare(j1, j2) == 0
}
