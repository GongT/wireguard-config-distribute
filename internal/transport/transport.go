package transport

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type Transport struct {
	publicListenPort  uint16
	realWireguardPort uint16
	realWireguardAddr *net.UDPAddr

	publicListen    *net.UDPConn
	mainLoopRunning bool

	natsMu sync.RWMutex
	nats   natRecords

	dispose       func()
	programIsQuit bool

	dumpTimer bool
}

func NewTransport() *Transport {
	return &Transport{}
}

func (t *Transport) Quit() {
	t.programIsQuit = true
	if t.Enabled() {
		tools.Debug("Stop software port forwarding")
		t.Stop()
	}
}

func (t *Transport) Stop() {
	if !t.Enabled() {
		// tools.Debug("transport status error: duplicate stop()")
		return
	}
	t.dispose()
	t.publicListen.Close()
	t.publicListen = nil
	t.closeNatMap()
}

func (t *Transport) Enabled() bool {
	return t.publicListen != nil
}

func (t *Transport) Start(port uint16) (uint16, error) {
	if t.Enabled() {
		if port != t.publicListenPort {
			tools.Die("transport status error: duplicate start() with different port!")
		}
		return t.realWireguardPort, nil
	}

	if t.realWireguardPort == 0 {
		if p, err := getFree(); err == nil {
			t.realWireguardPort = p
			t.realWireguardAddr = parse("127.0.0.1", p)
		} else {
			return 0, err
		}
	}

	addr, _ := net.ResolveUDPAddr("udp", format("0.0.0.0", port))
	if lis, err := net.ListenUDP("udp", addr); err != nil {
		return 0, err
	} else {
		if err := lis.SetReadBuffer(2000); err != nil {
			tools.Error("warn: set read buffer failed: %v", err)
		}
		t.publicListenPort = port
		t.publicListen = lis
	}

	t.initNatMap()
	t.goHandle()

	return t.realWireguardPort, nil
}

func (t *Transport) Sdump() (ret string) {
	t.natsMu.RLock()
	defer t.natsMu.RUnlock()
	ret += "Software port forwarding: "
	if t.Enabled() {
		ret += "Enabled\n"
	} else {
		ret += "Disabled\n"
		return
	}
	ret += fmt.Sprintf("    publicListenPort = %v", t.publicListenPort)
	if !t.mainLoopRunning {
		ret += " (Not Running!)"
	}
	ret += "\n"

	ret += fmt.Sprintf("    realWireguardPort = %v\n", t.realWireguardPort)

	for index, nat := range t.nats {
		ret += fmt.Sprintf("    [%v]: %v\n", index, nat.sdump())
		if nat.knownSessionId {
			ret += fmt.Sprintf("         Session Id: %v\n", nat.sessionId)
		} else {
			ret += fmt.Sprintf("         Last Alive: %v\n", nat.lastAlive.Format(time.RFC3339))
		}
	}
	ret += fmt.Sprintf("Total %v active forward.", len(t.nats))
	return
}
func (t *Transport) Dump() {
	if tools.IsDevelopmennt() && !t.dumpTimer {
		t.dumpTimer = true
		go func() {
			time.Sleep(5 * time.Second)
			t.dumpTimer = false

			tools.Debug(t.Sdump())
		}()
	}
}
