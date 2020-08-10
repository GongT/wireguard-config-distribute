package grpcImplements

import (
	"errors"
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (s *Implements) Start(req *protocol.IdReportingRequest, sender protocol.WireguardApi_StartServer) error {
	sid := types.DeSerializeSidType(req.SessionId)

	tools.Error("[%v] attached sender", sid)

	if !s.peersManager.AttachSender(sid, &sender) {
		return errors.New(fmt.Sprintf("Failed find client [%v] in registry", sid))
	}

	tools.Error("[%v] start loop", sid)
	<-sender.Context().Done()
	s.peersManager.Delete(sid)
	tools.Error("[%v] start return", sid)

	return nil
}
