package server

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (stat *ServerStatus) Greeting(request *protocol.ClientInfoRequest) (*protocol.ClientInfoResponse, error) {
	return stat.rpc.Greeting(stat.context, request)
}

func (stat *ServerStatus) Start() (protocol.WireguardApi_StartClient, error) {
	return stat.rpc.Start(stat.context, tools.EmptyPb)
}

func (stat *ServerStatus) Close() error {
	_, err := stat.rpc.Close(stat.context, tools.EmptyPb)
	return err
}

func (stat *ServerStatus) KeepAlive() (*protocol.KeepAliveStatus, error) {
	return stat.rpc.KeepAlive(stat.context, tools.EmptyPb)
}

func (stat *ServerStatus) NewGroup(request *protocol.NewGroupRequest) error {
	_, err := stat.rpc.NewGroup(stat.context, request)
	return err
}

/*
func (stat *ServerStatus) Register(request *protocol.RegisterRequest) (*protocol.RegisterResponse, error) {
	return stat.rpc.Register(stat.context, request)
}
func (stat *ServerStatus) Unregister(request *protocol.UnregisterRequest) error {
	_, err := stat.rpc.Unregister(stat.context, request)
	return err
}
*/
