package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (srv *Implements) KeepAlive(_ context.Context, request *protocol.IdReportingRequest) (*protocol.KeepAliveStatus, error) {
	defer tools.TimeMeasure("Grpc:KeepAlive")()

	sid := types.DeSerializeSidType(request.GetSessionId())
	succ := srv.peersManager.UpdateKeepAlive(sid)
	return &protocol.KeepAliveStatus{Success: succ}, nil
}
