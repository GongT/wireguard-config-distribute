package peerStatus

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/constants"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (peers *PeersManager) CleanupTimeoutPeers() {
	defer peers.m.Lock("CleanupTimeoutPeers")()

	tools.Debug(" ~ timer, do cleanup")
	expired := time.Now().Add(-5 * constants.KEEY_ALIVE_SECONDS)
	for cid, peer := range peers.list {
		if peer.lastKeepAlive.Before(expired) {
			tools.Error("[%v] peer exired", cid.Serialize())
			peers._delete(cid)
			peers.onChange.Write(cid)
		}
	}
}

func (peers *PeersManager) UpdateKeepAlive(cid types.SidType) bool {
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
