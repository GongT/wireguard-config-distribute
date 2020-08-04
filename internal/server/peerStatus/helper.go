package peerStatus

import (
	"bytes"
	"encoding/json"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func exactSame(a lPeerData, b lPeerData) bool {
	j1, err1 := json.Marshal(a)
	j2, err2 := json.Marshal(b)

	if err1 != nil || err2 != nil {
		return false
	}

	return bytes.Compare(j1, j2) == 0
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

func (peers *PeerStatus) sendSnapshot(peer lPeerData) {
	tools.Debug("[%v|%v] ~ send peers", peer.sessionId, peer.MachineId)
	err := (*peer.sender).Send(peers.createAllView(peer))
	if err == nil {
		tools.Debug("[%v|%v] ~ send peers ok", peer.sessionId, peer.MachineId)
	} else {
		tools.Debug("[%v|%v] ~ send peers failed: %s", peer.sessionId, peer.MachineId, err.Error())
	}
}
