package server

import (
	"fmt"
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements"
	"github.com/gongt/wireguard-config-distribute/internal/server/serverAuth"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
)

type listenOptions interface {
	GetListenPath() string
	GetListenPort() uint16

	GetPassword() string
}

type serverStateHolder struct {
	creds   *credentials.TransportCredentials
	service *grpcImplements.Implements

	listenUnix bool
	_listen    func()

	listener net.Listener
	grpc     *grpc.Server

	isQuit bool
}

func NewServer(options listenOptions, creds *credentials.TransportCredentials, serviceImpl *grpcImplements.Implements) *serverStateHolder {
	ret := &serverStateHolder{
		creds:   creds,
		service: serviceImpl,

		isQuit: false,
	}
	if len(options.GetListenPath()) > 0 {
		ret.creds = nil

		ret._listen = func() {
			ret.listenUnix = true
			ret.listener = listenUnix(options.GetListenPath())
		}
	} else if options.GetListenPort() > 0 {
		ret._listen = func() {
			ret.listenUnix = false
			ret.listener = listenTCP(options.GetListenPort())
		}
	} else {
		tools.Die("invalid config: no path or port to listen")
	}

	grpc.EnableTracing = tools.IsDevelopmennt()

	serverOptions := []grpc.ServerOption{}
	if ret.creds != nil {
		serverOptions = append(serverOptions, grpc.Creds(*ret.creds))
	}
	if pwd := options.GetPassword(); len(pwd) > 0 {
		handler := serverAuth.CreatePasswordCheck(pwd)
		serverOptions = append(serverOptions, grpc.StreamInterceptor(handler.Stream), grpc.UnaryInterceptor(handler.Unary))
	}
	ret.grpc = grpc.NewServer(serverOptions...)

	protocol.RegisterWireguardApiServer(ret.grpc, serviceImpl)

	if tools.IsDevelopmennt() {
		service.RegisterChannelzServiceToServer(ret.grpc)
	}

	return ret
}

func (srv *serverStateHolder) Listen(options listenOptions) {
	if srv.listener != nil {
		tools.Die("Program state error: grpc already listening")
	}

	srv._listen()
	srv._listen = nil

	listenAddr := srv.listener.Addr().String()
	go func() {
		fmt.Printf("Server listen on: %s\n", listenAddr)
		err := srv.grpc.Serve(srv.listener)
		if err == nil {
			tools.Error("Server listen successfully complete")
		} else {
			tools.Die("Server listen unexpected return: %s", err.Error())
		}
	}()
}

func (srv *serverStateHolder) Stop() {
	if srv.isQuit {
		tools.Error("Duplicate call to serverStateHolder.Stop()")
		return
	}
	srv.isQuit = true

	srv.grpc.Stop()
}

func (srv *serverStateHolder) IsSecure() bool {
	return srv.creds != nil
}
