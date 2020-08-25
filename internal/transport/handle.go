package transport

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type natRecord struct {
	knownSessionId bool
	sessionId      uint64    // one of
	lastAlive      time.Time // one of

	publicSocket      *net.UDPConn
	wgCommunicatePort uint16
	wgCommunicateConn *net.UDPConn

	remoteAddr *net.UDPAddr

	eventLoopRunning bool
	mu               sync.RWMutex
}

type natRecords = []*natRecord

func (nat *natRecord) sdump() string {
	var state string
	if nat.eventLoopRunning {
		state = "⚯"
	} else {
		state = "✂"
	}
	var dir string
	if nat.knownSessionId {
		dir = "->"
	} else {
		dir = "<-"
	}

	return fmt.Sprintf("realWg %v %v %v pubPort %v %v", dir, nat.wgCommunicatePort, state, dir, nat.remoteAddr.String())
}
