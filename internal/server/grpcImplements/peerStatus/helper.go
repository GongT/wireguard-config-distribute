package peerStatus

import (
	"bytes"
	"encoding/json"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func exactSame(a *PeerData, b *PeerData) bool {
	j1, err1 := json.Marshal(a)
	j2, err2 := json.Marshal(b)

	if err1 != nil || err2 != nil {
		return false
	}

	return bytes.Compare(j1, j2) == 0
}

func (peers *PeersManager) createSessionId(peer *PeerData) types.SidType {
	s := peer.CreateId()
	if sid, ok := peers.guidMap[s]; ok {
		peer.sessionId = sid
		return sid
	}

	peers.guid += 1
	peers.guidMap[s] = peers.guid
	peer.sessionId = peers.guid
	return peers.guid
}

func (peers *PeersManager) sendSnapshot(peer *PeerData) {
	tools.Debug("[%v] ~ send peers -> %s", peer.sessionId, peer.Title)
	list := peers.mapper[peer.VpnId]
	err := (*peer.sender).Send(list.generateAllView(peer))
	if err != nil {
		tools.Debug("[%v|%v] ~ send peers failed: %s", peer.sessionId, peer.MachineId, err.Error())
	}
}
