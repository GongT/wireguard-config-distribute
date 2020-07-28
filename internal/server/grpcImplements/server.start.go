package grpcImplements

import (
	"errors"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *Implements) Start(req *protocol.IdReportingRequest, sender protocol.WireguardApi_StartServer) error {
	sid := req.SessionId

	tools.Error("[%d] attached sender", sid)

	if !s.peersManager.AttachSender(sid, &sender) {
		return errors.New("Failed find client [" + strconv.FormatUint(sid, 10) + "] in registry")
	}

	tools.Error("[%d] start loop", sid)
	for {
		if sender.RecvMsg(nil) == nil {
			tools.Error("[%d] receive return nil", sid)
			break
		}
	}
	s.peersManager.Delete(sid)

	tools.Error("[%d] start return", sid)

	return nil
}
