package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
)

func (srv *Implements) KeepAlive(_ context.Context, request *protocol.IdReportingRequest) (*protocol.KeepAliveStatus, error) {
	sid := request.GetMachineId()
	succ := srv.peersManager.UpdateKeepAlive(sid)
	return &protocol.KeepAliveStatus{Success: succ}, nil
}
