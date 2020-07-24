package server

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
)

func (s *serverImplement) Greeting(context.Context, *protocol.ClientInfoRequest) (*protocol.ClientInfoResponse, error) {
	return nil, nil
}
