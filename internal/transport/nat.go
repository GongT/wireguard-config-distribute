package transport

import (
	"net"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (t *Transport) getNatFromOutside(remote *net.UDPAddr) *natRecord {
	t.natsMu.RLock()
	defer t.natsMu.RUnlock()

	for _, nat := range t.nats {
		if nat.remoteAddr.IP.Equal(remote.IP) && nat.remoteAddr.Port == remote.Port {
			if !nat.knownSessionId {
				nat.lastAlive = time.Now()
			}

			return nat
		}
	}
	return nil
}

func (t *Transport) initNatMap() {
	t.natsMu.Lock()
	defer t.natsMu.Unlock()

	l := len(t.nats)
	if l < 10 {
		l = 10
	}
	t.nats = make(natRecords, 0, l)
}

func (t *Transport) closeNatMap() {
	t.natsMu.Lock()
	defer t.natsMu.Unlock()

	for _, item := range t.nats {
		if err := item.stopOutgoingEventLoop(nil); err != nil {
			tools.Error("warn: failed close nat listen (%v->%v): %v", item.wgCommunicatePort, item.remoteAddr.String(), err)
		}
	}
	t.nats = nil
}

func maybeChangeRemote(nat *natRecord, targetIp string, targetPort uint16) {
	nat.mu.RLock()
	oldAddr := nat.remoteAddr.IP.String()
	oldPort := nat.remoteAddr.Port
	nat.mu.RUnlock()

	if oldAddr != targetIp || oldPort != int(targetPort) {
		tools.Debug(" ! peer ip updated: %v:%v -> %v:%v", oldAddr, oldPort, targetIp, targetPort)
		nat.mu.Lock()
		nat.remoteAddr = parse(targetIp, targetPort)
		nat.mu.Unlock()
	}
}
