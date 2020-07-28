package client

import (
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *clientStateHolder) Quit() {
	if s.isQuit {
		tools.Error("Duplicate call to Client.quit()")
		return
	}
	s.isQuit = true

	s.server.Disconnect(s.isRunning, s.MachineId)

	s.quitChan <- true
}
