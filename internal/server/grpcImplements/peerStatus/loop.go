package peerStatus

import (
	"fmt"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/systemd"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (peers *PeersManager) StopHandleChange() {
	peers.onChange.Close()
}

func (peers *PeersManager) StartHandleChange() {
	for changeCid := range peers.onChange.Read() {
		unlock := peers.m.Lock(fmt.Sprintf("StartHandleChange[%v]", changeCid))

		done := tools.TimeMeasure(fmt.Sprintf("peers::onChange[%v]::sendPeers", changeCid.Serialize()))

		len := len(peers.list)
		for cid, peer := range peers.list {
			if cid == changeCid || peer.sender == nil {
				continue
			}

			peers.sendSnapshot(peer)
		}

		unlock()

		systemd.UpdateState("peers(" + strconv.FormatInt(int64(len), 10) + ")")

		done()
	}
}
