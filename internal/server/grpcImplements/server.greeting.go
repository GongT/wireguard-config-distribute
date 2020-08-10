package grpcImplements

import (
	"context"
	"errors"
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/constants"
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements/peerStatus"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
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

	vpnName := request.GetGroupName()
	vpn, ok := s.vpnManager.GetLocked(vpnName)
	if !ok {
		return nil, errors.New("VPN group not exists: " + vpnName)
	}
	defer vpn.Release()
	fmt.Printf("   * VPN: %v\n", vpnName)

	networkGroup := request.GetNetwork().GetNetworkId()

	subnet := vpn.Subnet()
	if subnet == 0 {
		return nil, errors.New("VPN group config error: subnet is 0")
	}
	fmt.Printf("   * subnet: %v\n", subnet)

	allocIp := vpn.AllocateIp(request.GetHostname(), request.GetRequestVpnIp())
	if len(allocIp) == 0 {
		return nil, errors.New("Can not alloc ip address")
	}
	fmt.Printf("   * allocated ip address: %v\n", allocIp)

	externalIps := request.GetNetwork().GetExternalIp()
	if len(externalIps) == 0 {
		if request.GetNetwork().GetExternalEnabled() {
			externalIps = append(externalIps, remoteIp)
		}
	}
	fmt.Printf("   * external ips: %v\n", externalIps)

	pubKey, priKey, err := wireguard.GenerateKeyPair()
	if err != nil {
		return nil, errors.New("Failed generate wireguard keys: " + err.Error())
	}
	fmt.Printf("   * wireguard public: %v\n", pubKey)

	clientId := request.GetMachineId()
	if len(clientId) == 0 {
		clientId = networkGroup + "::" + request.GetHostname()
	}
	fmt.Printf("   * client id: %v\n", clientId)

	sessionId := s.peersManager.Add(&peerStatus.PeerData{
		MachineId:    clientId,
		VpnId:        types.DeSerializeVpnIdType(vpnName),
		Title:        request.GetTitle() + " [AT] " + request.GetNetwork().GetNetworkId(),
		Hostname:     request.GetHostname(),
		PublicKey:    pubKey,
		VpnIp:        allocIp,
		MTU:          request.GetNetwork().GetMTU(),
		Hosts:        request.GetServices(),
		NetworkId:    networkGroup,
		ExternalIp:   externalIps,
		ExternalPort: port(request.GetNetwork().GetExternalPort()),
		InternalIp:   request.GetNetwork().GetInternalIp(),
		InternalPort: port(request.GetNetwork().GetInternalPort()),
	})

	fmt.Printf("   * new session id: %v\n", sessionId)

	return &protocol.ClientInfoResponse{
		SessionId:  sessionId.Serialize(),
		MachineId:  clientId,
		PublicIp:   remoteIp,
		OfferIp:    allocIp,
		PrivateKey: priKey,
		Subnet:     uint32(subnet),
	}, nil
}

func port(n uint32) uint32 {
	if n == 0 {
		return constants.DEFAULT_PORT_NUMBER
	} else {
		return n
	}
}
