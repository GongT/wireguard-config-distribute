package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverImplement) Close(_ context.Context, request *protocol.IdReportingRequest) (*emptypb.Empty, error) {
	tools.Error("Call to Close (from %v)", request.GetSessionId())
	return tools.EmptyPb, nil
}
