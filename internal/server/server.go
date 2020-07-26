package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	channelzWebDebug "github.com/rantav/go-grpc-channelz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
)

type listenOptions interface {
	GetListenPath() string
	GetListenPort() uint16
}

type serverStateHolder struct {
	creds  *credentials.TransportCredentials
	srvice protocol.WireguardApiServer

	listenUnix bool
	_listen    func()

	listener      net.Listener
	debugListener net.Listener
	grpc          *grpc.Server
}

func NewServer(options listenOptions, creds *credentials.TransportCredentials, srv protocol.WireguardApiServer) *serverStateHolder {
	ret := &serverStateHolder{
		creds:  creds,
		srvice: srv,
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

	if ret.creds == nil {
		ret.grpc = grpc.NewServer()
	} else {
		ret.grpc = grpc.NewServer(grpc.Creds(*ret.creds))
	}

	protocol.RegisterWireguardApiServer(ret.grpc, srv)
	service.RegisterChannelzServiceToServer(ret.grpc)

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

	if tools.IsDevelopmennt() {
		http.Handle("/", channelzWebDebug.CreateHandler("/", listenAddr))
		srv.debugListener = listenTCP(options.GetListenPort() + 1)
		go func() {
			tools.Error("Debug listen port: %v", srv.debugListener.Addr().String())
			http.Serve(srv.debugListener, nil)
			tools.Error("Debug listen successfully complete")
		}()
	}
}

func (srv *serverStateHolder) Stop() {
	srv.grpc.Stop()
	srv.debugListener.Close()
}

func (srv *serverStateHolder) IsSecure() bool {
	return srv.creds != nil
}
