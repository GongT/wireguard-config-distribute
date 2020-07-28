package peerStatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/server/peerStatus/asyncChan"
	"github.com/gongt/wireguard-config-distribute/internal/server/peerStatus/debugLocker"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type lPeerData *PeerData

type PeerData struct {
	MachineId string
	Title     string
	Hostname  string
	PublicKey string
	VpnIp     string
	KeepAlive uint32
	MTU       uint32
	Hosts     []string

	NetworkId    string
	ExternalIp   []string
	ExternalPort uint32
	InternalIp   []string
	InternalPort uint32

	sender        *protocol.WireguardApi_StartServer
	lastKeepAlive time.Time
}

type PeerStatus struct {
	list     map[string]lPeerData
	m        debugLocker.MyLocker
	onChange *asyncChan.AsyncChan
}

func NewPeerStatus() *PeerStatus {
	return &PeerStatus{
		list:     make(map[string]lPeerData),
		onChange: asyncChan.NewChan(),
		m:        debugLocker.NewMutex(),
	}
}

func (peers *PeerStatus) StopHandleChange() {
	peers.onChange.Close()
}

func (peers *PeerStatus) AttachSender(cid string, sender *protocol.WireguardApi_StartServer) bool {
	defer peers.m.Lock(fmt.Sprintf("AttachSender[%s]", cid))()

	peer, exists := peers.list[cid]

	if !exists {
		return false
	}

	peer.sender = sender

	peers.sendSnapshot(peer)

	return true
}

func (peers *PeerStatus) sendSnapshot(peer lPeerData) {
	tools.Debug("[%s] ~ send peers", peer.MachineId)
	err := (*peer.sender).Send(peers.createAllView(peer))
	if err == nil {
		tools.Debug("[%s] ~ send peers ok", peer.MachineId)
	} else {
		tools.Debug("[%s] ~ send peers failed: %s", peer.MachineId, err.Error())
	}
}

func (peers *PeerStatus) StartHandleChange() {
	for changeCid := range peers.onChange.Read() {
		unlock := peers.m.Lock(fmt.Sprintf("StartHandleChange[%s]", changeCid))

		for cid, peer := range peers.list {
			if cid == changeCid || peer.sender == nil {
				continue
			}

			peers.sendSnapshot(peer)
		}

		unlock()
	}
}

func (peers *PeerStatus) DoCleanup() {
	defer peers.m.Lock("DoCleanup")()

	tools.Debug(" ~ timer, do cleanup")
	expired := time.Now().Add(-1 * time.Minute)
	for cid, peer := range peers.list {
		if peer.lastKeepAlive.Before(expired) {
			tools.Error("[%s] peer exired", peer.MachineId)
			delete(peers.list, cid)
			peers.onChange.Write(cid)
		}
	}
}
func (peers *PeerStatus) UpdateKeepAlive(cid string) bool {
	defer peers.m.Lock(fmt.Sprintf("UpdateKeepAlive[%s]", cid))()

	if peer, exists := peers.list[cid]; exists {
		tools.Debug("[%s] ~ keep alive", peer.MachineId)
		peer.lastKeepAlive = time.Now()
		return true
	} else {
		tools.Error("[%s] ! keep alive not exists peer", cid)
		return false
	}
}
func (peers *PeerStatus) Delete(cid string) {
	defer peers.m.Lock(fmt.Sprintf("Delete[%s]", cid))()

	_, exists := peers.list[cid]

	if !exists {
		tools.Error("[%s] ! delete not exists peer", cid)
		return
	}

	tools.Debug("[%s] ~ delete peer", cid)

	delete(peers.list, cid)
	peers.onChange.Write(cid)
}

func (peers *PeerStatus) Add(peer lPeerData) {
	defer peers.m.Lock(fmt.Sprintf("Add[%s]", peer.MachineId))()

	old, exists := peers.list[peer.MachineId]

	if exists {
		if exactSame(old, peer) {
			tools.Debug(" ~ add peer has exact same one")
			return
		}
		tools.Debug(" ~ replace peer")

		peer.lastKeepAlive = old.lastKeepAlive
		peer.sender = old.sender
	} else {
		tools.Debug(" ~ add new peer")
	}

	peers.list[peer.MachineId] = peer
	peers.onChange.Write(peer.MachineId)
}

func (peers *PeerStatus) createAllView(viewer lPeerData) *protocol.Peers {
	hosts := make(map[string]string)
	list := make([]*protocol.Peers_Peer, 0, len(peers.list)-1)

	for cid, peer := range peers.list {
		if viewer.MachineId == cid {
			continue
		}

		list = append(list, peers.createOneView(viewer, peer))

		for _, host := range peer.Hosts {
			hosts[host+"."+peer.Hostname] = peer.VpnIp
		}
	}

	return &protocol.Peers{
		List:  list,
		Hosts: hosts,
	}
}

func (peers *PeerStatus) createOneView(viewer lPeerData, peer lPeerData) *protocol.Peers_Peer {
	port := peer.ExternalPort
	ip := peer.ExternalIp
	if viewer.NetworkId == peer.NetworkId && len(viewer.NetworkId) > 0 {
		port = peer.InternalPort
		ip = peer.InternalIp
	}

	p := protocol.Peers_Peer{
		MachineId: peer.MachineId,
		Title:     peer.Title,
		Hostname:  peer.Hostname,
		Peer: &protocol.Peers_ConnectionTarget{
			PublicKey: peer.PublicKey,
			Address:   ip,
			Port:      port,
			VpnIp:     peer.VpnIp,
			KeepAlive: peer.KeepAlive,
			MTU:       peer.MTU,
		},
	}

	return &p
}

func exactSame(a lPeerData, b lPeerData) bool {
	j1, err1 := json.Marshal(a)
	j2, err2 := json.Marshal(b)

	if err1 != nil || err2 != nil {
		return false
	}

	return bytes.Compare(j1, j2) == 0
}
