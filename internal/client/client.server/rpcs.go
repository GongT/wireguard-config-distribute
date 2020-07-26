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

func (stat *ServerStatus) Close(id uint64) error {
	_, err := stat.rpc.Close(stat.context, &protocol.IdReportingRequest{SessionId: id})
	return err
}

func (stat *ServerStatus) KeepAlive(id uint64) (*protocol.KeepAliveStatus, error) {
	return stat.rpc.KeepAlive(stat.context, &protocol.IdReportingRequest{SessionId: id})
}

func (stat *ServerStatus) NewGroup(request *protocol.NewGroupRequest) error {
	_, err := stat.rpc.NewGroup(stat.context, request)
	return err
}

func (stat *ServerStatus) GetSelfSignedCertFile(request *protocol.GetCertFileRequest) (*protocol.GetCertFileResponse, error) {
	return stat.rpc.GetSelfSignedCertFile(stat.context, request)
}
