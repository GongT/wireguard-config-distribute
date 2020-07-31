package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (srv *Implements) KeepAlive(_ context.Context, request *protocol.IdReportingRequest) (*protocol.KeepAliveStatus, error) {
	sid := types.DeSerialize(request.GetSessionId())
	succ := srv.peersManager.UpdateKeepAlive(sid)
	return &protocol.KeepAliveStatus{Success: succ}, nil
}
