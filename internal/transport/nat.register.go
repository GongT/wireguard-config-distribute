package transport

import (
	"github.com/gongt/wireguard-config-distribute/internal/client/clientType"
)

func (t *Transport) ModifyPeers(list clientType.PeerDataList) {
	if !t.Enabled() {
		return
	}
	for _, client := range list {
		peer := client.GetPeer()
		if peer.GetSameNetwork() {
			continue
		}
		address := client.GetSelectedAddress()
		if len(address) == 0 {
			continue
		}

		p := uint16(peer.GetPort())

		sid := client.GetSessionId()
		nat := t.FindById(sid)
		if nat == nil {
			nat = createNatOutgoing(address, p)
			nat.knownSessionId = true
			nat.sessionId = sid
			nat.publicSocket = t.publicListen

			t.natsMu.Lock()
			t.nats = append(t.nats, nat)
			t.natsMu.Unlock()

			nat.goOutgoingEventLoop()
		} else {
			maybeChangeRemote(nat, address, p)
		}
		client.ChangeTo("127.0.0.1", nat.wgCommunicatePort)
	}

	t.Dump()
}

func (t *Transport) FindById(sid uint64) *natRecord {
	t.natsMu.RLock()
	defer t.natsMu.RUnlock()

	for _, search := range t.nats {
		if search.knownSessionId && search.sessionId == sid {
			return search
		}
	}
	return nil
}
