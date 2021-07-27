package peerStatus

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type LockedPeerData struct {
	*PeerData
	Unlock func()
}

func (peers *PeersManager) GetLocked(sid types.SidType) *LockedPeerData {
	unlock := peers.m.Lock(fmt.Sprintf("GetLocked[%v]", sid))

	if _, exists := peers.list[sid]; !exists {
		tools.Error("grpc:UpdateClientInfo() fail: %v not exists in:", sid)
		for k, p := range peers.list {
			tools.Error("  %v: %s", k, p.MachineId)
		}
		unlock()
		return nil
	}

	return &LockedPeerData{
		PeerData: peers.list[sid],
		Unlock: func() {
			unlock()
			peers.onChange.Write(sid)
		},
	}
}

func (peers *PeersManager) AttachSender(sid types.SidType, sender *protocol.WireguardApi_StartServer) bool {
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

	done := tools.TimeMeasure(fmt.Sprintf("peers::AttachSender[%v]::sendSnapshot", sid.Serialize()))
	peers.sendSnapshot(peer)
	done()

	return true
}

func (peers *PeersManager) _delete(cid types.SidType) error {
	old, exists := peers.list[cid]
	if !exists {
		return fmt.Errorf("[%v] ! delete not exists peer", cid)
	}
	tools.Debug("[%v] ~ delete peer", cid)
	vpn := old.VpnId
	delete(peers.mapper[vpn], cid)
	delete(peers.list, cid)
	if len(peers.mapper[vpn]) == 0 {
		tools.Debug(" ~ all peer deleted. (%v)", vpn.Serialize())
		delete(peers.mapper, vpn)
	}
	return nil
}

func (peers *PeersManager) Delete(cid types.SidType) {
	defer peers.m.Lock(fmt.Sprintf("Delete[%v]", cid))()
	peers._delete(cid)
	peers.onChange.Write(cid)
}

func (peers *PeersManager) Add(peer *PeerData) (sid types.SidType) {
	defer peers.m.Lock(fmt.Sprintf("Add[%v]", peer.MachineId))()

	sid = peers.createSessionId(peer)
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

	if _, ok := peers.mapper[peer.VpnId]; !ok {
		peers.mapper[peer.VpnId] = make(vpnPeersMap)
	}
	peers.mapper[peer.VpnId][sid] = peer

	return
}
