package peerStatus

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

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
