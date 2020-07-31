package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implements) Close(_ context.Context, request *protocol.IdReportingRequest) (*emptypb.Empty, error) {
	s.peersManager.Delete(types.DeSerialize(request.GetSessionId()))
	return tools.EmptyPb, nil
}
