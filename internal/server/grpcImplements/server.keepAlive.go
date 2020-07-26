package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *serverImplement) KeepAlive(_ context.Context, request *protocol.IdReportingRequest) (*protocol.KeepAliveStatus, error) {
	tools.Error("Call to KeepAlive (from %v)", request.GetSessionId())
	return &protocol.KeepAliveStatus{Success: true}, nil
}
