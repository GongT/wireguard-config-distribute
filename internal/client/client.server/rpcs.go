package server

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
)

func (stat *ServerStatus) Greeting(request *protocol.ClientInfoRequest) (*protocol.ClientInfoResponse, error) {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	return stat.rpc.Greeting(ctx, request)
}

func (stat *ServerStatus) Start(id uint64) (<-chan *protocol.Peers, error) {
	cctx, cancel := context.WithCancel(stat.context)
	stream, err := stat.rpc.Start(cctx, &protocol.IdReportingRequest{SessionId: id})
	if err != nil {
		return nil, err
	}

	ch := make(chan *protocol.Peers)

	go func() {
		for {
			peers, err := stream.Recv()
			if err != nil {
				break
			}
			ch <- peers
		}
		cancel()
		close(ch)
	}()

	return ch, nil
}

func (stat *ServerStatus) Close(id uint64) error {
	ctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	_, err := stat.rpc.Close(ctx, &protocol.IdReportingRequest{SessionId: id})
	return err
}

func (stat *ServerStatus) KeepAlive(id uint64) (*protocol.KeepAliveStatus, error) {
	cctx, cancel := context.WithCancel(stat.context)
	defer cancel()
	return stat.rpc.KeepAlive(cctx, &protocol.IdReportingRequest{SessionId: id})
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
