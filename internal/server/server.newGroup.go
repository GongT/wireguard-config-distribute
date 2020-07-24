package server

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverImplement) NewGroup(context.Context, *protocol.NewGroupRequest) (*emptypb.Empty, error) {
	return nil, nil
}
