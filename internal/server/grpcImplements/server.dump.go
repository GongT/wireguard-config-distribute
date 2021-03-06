package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implements) DumpStatus(_ context.Context, _ *emptypb.Empty) (*protocol.DumpResponse, error) {
	ret1 := s.vpnManager.Dump()

	ret2 := s.peersManager.Dump()
	return &protocol.DumpResponse{
		Text: ret1 + "\n\n" + ret2,
	}, nil
}
