package server

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverImplement) Start(*emptypb.Empty, protocol.WireguardApi_StartServer) error {
	return nil
}
