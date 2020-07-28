package grpcImplements

import (
	"errors"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *Implements) Start(req *protocol.IdReportingRequest, sender protocol.WireguardApi_StartServer) error {
	sid := req.MachineId

	tools.Error("[%s] attached sender", sid)

	if !s.peersManager.AttachSender(sid, &sender) {
		return errors.New("Failed find client [" + sid + "] in registry")
	}

	tools.Error("[%s] start loop", sid)
	for {
		if sender.RecvMsg(nil) == nil {
			tools.Error("[%s] receive return nil", sid)
			break
		}
	}
	s.peersManager.Delete(sid)

	tools.Error("[%s] start return", sid)

	return nil
}
