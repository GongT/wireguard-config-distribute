package grpcImplements

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implements) NewGroup(context.Context, *protocol.NewGroupRequest) (*emptypb.Empty, error) {
	return nil, nil
}
