package server

import (
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"google.golang.org/grpc"
)

type server struct {
	grpc *grpc.Server
}

type serverImplement struct {
}

func NewServer() (srv server) {
	srv = server{}

	grpcServer := grpc.NewServer()
	protocol.RegisterWireguardApiServer(grpcServer, &serverImplement{})

	srv.grpc = grpcServer

	return
}

func (s *server) ListenSocket(lis net.Listener) {
	s.grpc.Serve(lis)
}
