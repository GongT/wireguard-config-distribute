package transport

import (
	"net"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func createNatOutgoing(targetIp string, targetPort uint16) *natRecord {
	ret := natRecord{
		remoteAddr: parse(targetIp, targetPort),
	}

	// for wireguard init connect out
	lis, err := net.ListenUDP("udp", &net.UDPAddr{})

	if err != nil {
		tools.Die("failed listen UDP: %v", err)
	}
	if err := lis.SetReadBuffer(2000); err != nil {
		tools.Error("warn: set read buffer failed: %v", err)
	}

	ret.wgCommunicateConn = lis
	ret.wgCommunicatePort = uint16(lis.LocalAddr().(*net.UDPAddr).Port)

	return &ret
}

func (nat *natRecord) goOutgoingEventLoop() {
	nat.mu.Lock()
	defer nat.mu.Unlock()

	nat.eventLoopRunning = true
	go nat.eventLoopOutgoing()
}

func (nat *natRecord) eventLoopOutgoing() {
	buff := newBuffer()
	for {
		// 这里是从本机的wireguard收到的数据
		n, _, err := nat.wgCommunicateConn.ReadFromUDP(buff[:])

		if err != nil {
			if !isSocketClosed(err) {
				tools.Error("warn: forward %v:\n    local wireguard not listening: %v", nat.sdump(), err)
			}
			break
		}

		if dataDump {
			if nat.knownSessionId {
				tools.Debug("wg -> [%v bytes] -> %v", n, nat.sessionId)
			} else {
				tools.Debug("wg -> [%v bytes] -> %v", n, nat.remoteAddr.String())
			}
		}

		nat.mu.RLock()
		s := nat.publicSocket
		a := nat.remoteAddr
		if !nat.knownSessionId {
			nat.lastAlive = time.Now()
		}
		nat.mu.RUnlock()

		_, err = s.WriteToUDP(buff.encode(n), a)

		if err != nil {
			if isSocketClosed(err) {
				break
			} else {
				tools.Error("warn: write to wg failed: %v", err)
			}
		}
	}
	tools.Debug(" ! [%v] loop done", nat.sdump())
}

func (nat *natRecord) stopOutgoingEventLoop(err error) error {
	nat.mu.Lock()
	defer nat.mu.Unlock()

	if !nat.eventLoopRunning {
		return nil
	}

	if err != nil {
		tools.Error("warn: %v:\n    error when stop: %v", nat.sdump(), err)
	}

	nat.eventLoopRunning = false
	return nat.wgCommunicateConn.Close()
}
