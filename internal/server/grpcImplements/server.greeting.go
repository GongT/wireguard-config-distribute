package grpcImplements

import (
	"context"
	"errors"
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/constants"
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/server/peerStatus"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/wireguard"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func (s *Implements) Greeting(ctx context.Context, request *protocol.ClientInfoRequest) (*protocol.ClientInfoResponse, error) {
	remoteIp := tools.GetRemoteFromContext(ctx)
	if len(remoteIp) == 0 {
		return nil, errors.New("Failed find your ip")
	}

	authtype := "not auth"
	if p, _ := peer.FromContext(ctx); p.AuthInfo != nil {
		authtype = p.AuthInfo.AuthType()
	}

	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Printf("New Client Greeting: %s (%s)\n", remoteIp, authtype)
	for key, value := range md {
		fmt.Printf("   * %v: %v\n", key, value)
	}

	if !s.vpnManager.Exists(request.GetGroupName()) {
		return nil, errors.New("VPN group not exists: " + request.GetGroupName())
	}

	networkGroup := request.GetNetwork().GetNetworkId()

	allocIp := s.vpnManager.AllocateIp(request.GetGroupName(), request.GetHostname(), request.GetRequestVpnIp())
	if len(allocIp) == 0 {
		return nil, errors.New("Can not alloc ip address")
	}

	keepAlive := uint32(0)
	externalIps := request.GetNetwork().GetExternalIp()
	if len(externalIps) == 0 {
		if request.GetNetwork().GetExternalEnabled() {
			externalIps = append(externalIps, remoteIp)
			keepAlive = 25
		}
	}

	pubKey, priKey, err := wireguard.GenerateKeyPair()
	if err != nil {
		return nil, errors.New("Failed generate wireguard keys: " + err.Error())
	}

	clientId := request.GetMachineId()
	if len(clientId) == 0 {
		clientId = networkGroup + "::" + request.GetHostname()
	}

	s.peersManager.Add(&peerStatus.PeerData{
		MachineId:    clientId,
		Title:        request.GetTitle(),
		Hostname:     request.GetHostname(),
		PublicKey:    pubKey,
		VpnIp:        allocIp,
		KeepAlive:    keepAlive,
		MTU:          request.GetNetwork().GetMTU(),
		Hosts:        request.GetServices(),
		NetworkId:    networkGroup,
		ExternalIp:   externalIps,
		ExternalPort: port(request.GetNetwork().GetExternalPort()),
		InternalIp:   request.GetNetwork().GetInternalIp(),
		InternalPort: port(request.GetNetwork().GetInternalPort()),
	})

	return &protocol.ClientInfoResponse{
		MachineId:  clientId,
		PublicIp:   remoteIp,
		OfferIp:    allocIp,
		PrivateKey: priKey,
	}, nil
}

func port(n uint32) uint32 {
	if n == 0 {
		return constants.DEFAULT_PORT_NUMBER
	} else {
		return n
	}
}
