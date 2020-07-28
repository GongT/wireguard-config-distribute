package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implements) Close(_ context.Context, request *protocol.IdReportingRequest) (*emptypb.Empty, error) {
	s.peersManager.Delete(request.GetMachineId())
	return tools.EmptyPb, nil
}
