package peerStatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/asyncChannels"
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/server/peerStatus/debugLocker"
	"github.com/gongt/wireguard-config-distribute/internal/systemd"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type lPeerData *PeerData

type PeerData struct {
	sessionId types.SidType
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
	InternalIp   string
	InternalPort uint32

	sender        *protocol.WireguardApi_StartServer
	lastKeepAlive time.Time
}

type PeerStatus struct {
	list     map[types.SidType]lPeerData
	m        debugLocker.MyLocker
	onChange *asyncChannels.AsyncChanSidType

	guid    types.SidType
	guidMap map[string]types.SidType
}

func NewPeerStatus() *PeerStatus {
	return &PeerStatus{
		list:     make(map[types.SidType]lPeerData),
		onChange: asyncChannels.NewChanSidType(),
		m:        debugLocker.NewMutex(),

		guid:    0,
		guidMap: make(map[string]types.SidType),
	}
}

func (peers *PeerStatus) StopHandleChange() {
	peers.onChange.Close()
}

func (peers *PeerStatus) AttachSender(sid types.SidType, sender *protocol.WireguardApi_StartServer) bool {
	defer peers.m.Lock(fmt.Sprintf("AttachSender[%v]", sid))()

	peer, exists := peers.list[sid]

	if !exists {
		tools.Error("grpc:Start() fail: %v not exists in:", sid)
		for k, p := range peers.list {
			tools.Error("%v: %s", k, p.MachineId)
		}
		return false
	}

	peer.sender = sender

	peers.sendSnapshot(peer)

	return true
}

func (peers *PeerStatus) sendSnapshot(peer lPeerData) {
	tools.Debug("[%v|%v] ~ send peers", peer.sessionId, peer.MachineId)
	err := (*peer.sender).Send(peers.createAllView(peer))
	if err == nil {
		tools.Debug("[%v|%v] ~ send peers ok", peer.sessionId, peer.MachineId)
	} else {
		tools.Debug("[%v|%v] ~ send peers failed: %s", peer.sessionId, peer.MachineId, err.Error())
	}
}

func (peers *PeerStatus) StartHandleChange() {
	for changeCid := range peers.onChange.Read() {
		unlock := peers.m.Lock(fmt.Sprintf("StartHandleChange[%v]", changeCid))

		len := len(peers.list)
		for cid, peer := range peers.list {
			if cid == changeCid || peer.sender == nil {
				continue
			}

			peers.sendSnapshot(peer)
		}

		unlock()

		systemd.UpdateState("peers(" + strconv.FormatInt(int64(len), 10) + ")")
	}
}

func (peers *PeerStatus) DoCleanup() {
	defer peers.m.Lock("DoCleanup")()

	tools.Debug(" ~ timer, do cleanup")
	expired := time.Now().Add(-1 * time.Minute)
	for cid, peer := range peers.list {
		if peer.lastKeepAlive.Before(expired) {
			tools.Error("[%v] peer exired", peer.MachineId)
			delete(peers.list, cid)
			peers.onChange.Write(cid)
		}
	}
}
func (peers *PeerStatus) UpdateKeepAlive(cid types.SidType) bool {
	defer peers.m.Lock(fmt.Sprintf("UpdateKeepAlive[%v]", cid))()

	if peer, exists := peers.list[cid]; exists {
		tools.Debug("[%v|%v] ~ keep alive", cid, peer.MachineId)
		peer.lastKeepAlive = time.Now()
		return true
	} else {
		tools.Error("[%v] ! keep alive not exists peer", cid)
		return false
	}
}
func (peers *PeerStatus) Delete(cid types.SidType) {
	defer peers.m.Lock(fmt.Sprintf("Delete[%v]", cid))()

	_, exists := peers.list[cid]

	if !exists {
		tools.Error("[%v] ! delete not exists peer", cid)
		return
	}

	tools.Debug("[%v] ~ delete peer", cid)

	delete(peers.list, cid)
	peers.onChange.Write(cid)
}

func (peers *PeerStatus) createSessionId(networkId string, machineId string) types.SidType {
	s := networkId + "::" + machineId
	if sid, ok := peers.guidMap[s]; ok {
		return sid
	}

	peers.guid += 1
	peers.guidMap[s] = peers.guid
	return peers.guid
}

func (peers *PeerStatus) Add(peer lPeerData) (sid types.SidType) {
	defer peers.m.Lock(fmt.Sprintf("Add[%v]", peer.MachineId))()

	sid = peers.createSessionId(peer.NetworkId, peer.MachineId)
	peer.sessionId = sid
	old, exists := peers.list[sid]

	if exists {
		if exactSame(old, peer) {
			tools.Debug(" ~ add peer has exact same one")
			return
		}
		tools.Debug(" ~ replace peer")

		peer.lastKeepAlive = old.lastKeepAlive
		peer.sender = old.sender
	} else {
		peer.lastKeepAlive = time.Now()
		tools.Debug(" ~ add new peer")
	}

	peers.list[sid] = peer
	peers.onChange.Write(sid)

	return
}

func (peers *PeerStatus) createAllView(viewer lPeerData) *protocol.Peers {
	hosts := make(map[string]string)
	list := make([]*protocol.Peers_Peer, 0, len(peers.list)-1)

	for cid, peer := range peers.list {
		if viewer.sessionId == cid {
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
		ip = []string{peer.InternalIp}
	}

	p := protocol.Peers_Peer{
		SessionId: peer.sessionId.Serialize(),
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
