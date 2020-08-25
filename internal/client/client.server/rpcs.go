package server

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (stat *ServerStatus) RegisterClient(request *protocol.RegisterClientRequest) (*protocol.RegisterClientResponse, error) {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	return stat.rpc.RegisterClient(ctx, request)
}

func (stat *ServerStatus) UpdateClientInfo(request *protocol.ClientInfoRequest) (*protocol.ClientInfoResponse, error) {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	return stat.rpc.UpdateClientInfo(ctx, request)
}

func (stat *ServerStatus) Start(id uint64) (<-chan *protocol.Peers, error) {
	cctx, cancel := context.WithCancel(stat.context)
	stream, err := stat.rpc.Start(cctx, &protocol.IdReportingRequest{SessionId: id})
	if err != nil {
		cancel()
		return nil, err
	}

	ch := make(chan *protocol.Peers)

	go func() {
		for {
			peers, err := stream.Recv()
			if err != nil {
				tools.Debug(" ~ grpc:Start() disconnected: %s", err.Error())
				break
			}
			ch <- peers
		}
		cancel()
		close(ch)
	}()

	return ch, nil
}

func (stat *ServerStatus) Close(id types.SidType) error {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	_, err := stat.rpc.Close(ctx, &protocol.IdReportingRequest{SessionId: id.Serialize()})
	return err
}

func (stat *ServerStatus) KeepAlive(id types.SidType) (*protocol.KeepAliveStatus, error) {
	cctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	return stat.rpc.KeepAlive(cctx, &protocol.IdReportingRequest{SessionId: id.Serialize()})
}

func (stat *ServerStatus) NewGroup(request *protocol.NewGroupRequest) error {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	_, err := stat.rpc.NewGroup(ctx, request)
	return err
}

func (stat *ServerStatus) GetSelfSignedCertFile(request *protocol.GetCertFileRequest) (*protocol.GetCertFileResponse, error) {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	return stat.rpc.GetSelfSignedCertFile(ctx, request)
}

func (stat *ServerStatus) DumpStatus() (*protocol.DumpResponse, error) {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	return stat.rpc.DumpStatus(ctx, tools.EmptyPb)
}
