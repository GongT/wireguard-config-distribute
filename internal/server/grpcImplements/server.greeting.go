package grpcImplements

import (
	"context"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var _guid uint64 = 0

func guid() uint64 {
	_guid += 1
	return _guid
}

func (s *serverImplement) Greeting(ctx context.Context, _ *protocol.ClientInfoRequest) (*protocol.ClientInfoResponse, error) {
	remoteIp := tools.GetRemoteFromContext(ctx)
	if len(remoteIp) == 0 {
		return nil, errors.New("Failed find your ip")
	}

	authtype := "not auth"
	if p, _ := peer.FromContext(ctx); p.AuthInfo != nil {
		authtype = p.AuthInfo.AuthType()
	}

	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Printf("New Client Greeting: %s - %s: %s", remoteIp, authtype, spew.Sdump(md))

	clientId := guid()

	return &protocol.ClientInfoResponse{
		SessionId: clientId,
		OfferIp:   "1.2",
		PublicIp:  remoteIp,
	}, nil
}
