package server

import (
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc"
)

type serverGrpc struct {
	listen net.Listener
	grpc   *grpc.Server
}

type serverStateHolder struct {
	server
}

type serverImplement struct {
}

func NewServer(creds grpc.ServerOption) (srv serverStateHolder) {

	grpc.EnableTracing = tools.IsDevelopmennt()
	grpcServer := grpc.NewServer(creds)
	protocol.RegisterWireguardApiServer(grpcServer, &serverImplement{})

	srv.grpc = grpcServer

	return
}

type listenOptions interface {
	GetListenPath() string
	GetListenPort() uint16
}

func (s *serverStateHolder) Listen(options listenOptions) {
	if len(options.GetListenPath()) > 0 {
		go s.server.grpc.Serve(*listenUnix(options.GetListenPath()))
	} else if options.GetListenPort() > 0 {
		go s.server.grpc.Serve(*listenTCP(options.GetListenPort()))
	} else {
		tools.Die("invalid config: no path or port to listen")
	}
}

func (s *serverStateHolder) Stop() {
	s.grpc.Stop()
}
