package transport

import (
	"net"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func createNatIncomming(from *net.UDPAddr) *natRecord {
	ret := natRecord{
		remoteAddr: from,
	}

	// for wireguard init connect out
	conn, err := net.ListenUDP("udp", &net.UDPAddr{})

	if err != nil {
		tools.Die("failed listen UDP: %v", err)
	}
	if err := conn.SetReadBuffer(2000); err != nil {
		tools.Error("warn: set read buffer failed: %v", err)
	}

	ret.knownSessionId = false
	ret.lastAlive = time.Now()
	ret.wgCommunicateConn = conn
	ret.wgCommunicatePort = uint16(conn.LocalAddr().(*net.UDPAddr).Port)

	return &ret
}

func (t *Transport) handleErrorClosed() {
	t.natsMu.Lock()
	defer t.natsMu.Unlock()
	t.mainLoopRunning = false
}
func (t *Transport) goHandle() {
	t.natsMu.Lock()
	defer t.natsMu.Unlock()

	if t.mainLoopRunning {
		panic("go handle twice???")
	}
	t.mainLoopRunning = true
	go t.handlePublicIncomeConnect()
}

func (t *Transport) handlePublicIncomeConnect() {
	buff := newBuffer()
	for {
		// 这里是从远程发来的数据
		n, remote, err := t.publicListen.ReadFromUDP(buff[:])

		if err != nil {
			if isSocketClosed(err) {
				t.handleErrorClosed()
				return
			}
			tools.Die("failed read from public port: %v", err)
		}

		nat := t.getNatFromOutside(remote)
		if nat == nil {
			nat = createNatIncomming(remote)
			nat.publicSocket = t.publicListen

			t.natsMu.Lock()
			t.nats = append(t.nats, nat)
			t.natsMu.Unlock()

			nat.goOutgoingEventLoop()

			go t.Dump()
		}
		// wgAddr := parse("127.0.0.1", wgPort)
		n1, err := nat.wgCommunicateConn.WriteToUDP(buff.decode(n), t.realWireguardAddr)
		if err != nil {
			tools.Error("failed write to wireguard port: %v", err)
		} else if n != n1 {
			tools.Error("write bytes %v not equals to exptected %v", n1, n)
		}

		if dataDump {
			if nat.knownSessionId {
				tools.Debug("wg <- [%v bytes] <- %v", n, nat.sessionId)
			} else {
				tools.Debug("wg <- [%v bytes] <- %v", n, nat.remoteAddr.String())
			}
		}
	}
}
