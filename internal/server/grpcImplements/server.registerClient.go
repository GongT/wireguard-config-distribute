package grpcImplements

import (
	"context"
	"errors"
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements/peerStatus"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func returnError(err error) (*protocol.RegisterClientResponse, error) {
	werr := fmt.Errorf("greeting failed: %w", err)
	fmt.Printf("%v\n", werr)
	return nil, werr
}

func (s *Implements) RegisterClient(ctx context.Context, request *protocol.RegisterClientRequest) (*protocol.RegisterClientResponse, error) {
	remoteIp := tools.GetRemoteFromContext(ctx)
	if len(remoteIp) == 0 {
		return returnError(errors.New("failed find your ip"))
	}

	authtype := "not auth"
	if p, _ := peer.FromContext(ctx); p.AuthInfo != nil {
		authtype = p.AuthInfo.AuthType()
	}

	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Printf("New Client Greeting: %s (%s)\n", remoteIp, authtype)
	for key, value := range md {
		tools.Debug("   * %v: %v\n", key, value)
	}

	vpnName := types.DeSerializeVpnIdType(request.GetVpnGroup())
	vpn, ok := s.vpnManager.GetLocked(vpnName)
	if !ok {
		return returnError(errors.New("VPN group not exists: " + vpnName.Serialize()))
	}
	defer vpn.Release()
	fmt.Printf("   * VPN: %v\n", vpnName)

	networkGroup := request.GetLocalGroup()
	fmt.Printf("   * network group: %v\n", networkGroup)

	subnet := vpn.Subnet()
	if subnet == 0 {
		return returnError(errors.New("VPN group config error: subnet is 0"))
	}
	fmt.Printf("   * subnet: %v\n", subnet)

	allocIp := vpn.AllocateIp(request.GetHostname(), request.GetRequestVpnIp())
	if len(allocIp) == 0 {
		return returnError(errors.New("can not alloc ip address"))
	}
	fmt.Printf("   * allocated ip address: %v\n", allocIp)

	keys, err := vpn.AllocateKeyPair(request.GetHostname())
	if err != nil {
		return returnError(fmt.Errorf("failed generate wireguard keys: %w", err))
	}
	fmt.Printf("   * wireguard public: %v\n", keys.Public)

	clientId := request.GetMachineId()
	if len(clientId) == 0 {
		clientId = networkGroup + "::" + request.GetHostname()
	}
	fmt.Printf("   * client id: %v\n", clientId)

	sessionId := s.peersManager.Add(&peerStatus.PeerData{
		MachineId:   clientId,
		VpnId:       vpnName,
		Title:       request.GetTitle() + " [AT] " + request.GetLocalGroup(),
		Hostname:    request.GetHostname(),
		PublicKey:   keys.Public,
		VpnIp:       allocIp,
		WorkgroupId: networkGroup,
	})

	fmt.Printf("   * session id: %v\n", sessionId)

	return &protocol.RegisterClientResponse{
		SessionId:    sessionId.Serialize(),
		MachineId:    clientId,
		PublicIp:     remoteIp,
		OfferIp:      allocIp,
		PrivateKey:   keys.Private,
		Subnet:       uint32(subnet),
		EnableObfuse: vpn.GetObfuse(),
	}, nil
}
