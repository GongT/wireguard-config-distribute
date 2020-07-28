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
	SessionId uint64
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
	list     map[uint64]lPeerData
	m        debugLocker.MyLocker
	onChange *asyncChan.AsyncChan
}

func NewPeerStatus() *PeerStatus {
	return &PeerStatus{
		list:     make(map[uint64]lPeerData),
		onChange: asyncChan.NewChan(),
		m:        debugLocker.NewMutex(),
	}
}

func (peers *PeerStatus) StopHandleChange() {
	peers.onChange.Close()
}

func (peers *PeerStatus) AttachSender(cid uint64, sender *protocol.WireguardApi_StartServer) bool {
	defer peers.m.Lock(fmt.Sprintf("AttachSender[%d]", cid))()

	peer, exists := peers.list[cid]

	if !exists {
		return false
	}

	peer.sender = sender

	peers.sendSnapshot(peer)

	return true
}

func (peers *PeerStatus) sendSnapshot(peer lPeerData) {
	tools.Debug("[%d] ~ send peers", peer.SessionId)
	err := (*peer.sender).Send(peers.createAllView(peer))
	if err == nil {
		tools.Debug("[%d] ~ send peers ok", peer.SessionId)
	} else {
		tools.Debug("[%d] ~ send peers failed: %s", peer.SessionId, err.Error())
	}
}

func (peers *PeerStatus) StartHandleChange() {
	for changeCid := range peers.onChange.Read() {
		unlock := peers.m.Lock(fmt.Sprintf("StartHandleChange[%d]", changeCid))

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
			tools.Error("[%d] peer exired", peer.SessionId)
			delete(peers.list, cid)
			peers.onChange.Write(cid)
		}
	}
}
func (peers *PeerStatus) UpdateKeepAlive(cid uint64) bool {
	defer peers.m.Lock(fmt.Sprintf("UpdateKeepAlive[%d]", cid))()

	if peer, exists := peers.list[cid]; exists {
		tools.Debug("[%d] ~ keep alive", peer.SessionId)
		peer.lastKeepAlive = time.Now()
		return true
	} else {
		tools.Error("[%d] ! keep alive not exists peer", cid)
		return false
	}
}
func (peers *PeerStatus) Delete(cid uint64) {
	defer peers.m.Lock(fmt.Sprintf("Delete[%d]", cid))()

	_, exists := peers.list[cid]

	if !exists {
		tools.Error("[%d] ! delete not exists peer", cid)
		return
	}

	tools.Debug("[%d] ~ delete peer", cid)

	delete(peers.list, cid)
	peers.onChange.Write(cid)
}

func (peers *PeerStatus) Add(peer lPeerData) {
	defer peers.m.Lock(fmt.Sprintf("Add[%d]", peer.SessionId))()

	old, exists := peers.list[peer.SessionId]

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

	peers.list[peer.SessionId] = peer
	peers.onChange.Write(peer.SessionId)
}

func (peers *PeerStatus) createAllView(viewer lPeerData) *protocol.Peers {
	hosts := make(map[string]string)
	list := make([]*protocol.Peers_Peer, 0, len(peers.list)-1)

	for cid, peer := range peers.list {
		if viewer.SessionId == cid {
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
		SessionId: peer.SessionId,
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
