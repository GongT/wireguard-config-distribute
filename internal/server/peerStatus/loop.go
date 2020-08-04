package peerStatus

import (
	"fmt"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/systemd"
)

func (peers *PeerStatus) StopHandleChange() {
	peers.onChange.Close()
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
