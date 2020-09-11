package grpcImplements

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/constants"
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (s *Implements) UpdateClientInfo(ctx context.Context, request *protocol.ClientInfoRequest) (*protocol.ClientInfoResponse, error) {
	sessionId := types.SidType(request.GetSessionId())
	peerSession := s.peersManager.GetLocked(sessionId)

	if peerSession == nil {
		return nil, errors.New("not registered client!")
	}

	defer peerSession.Unlock()

	remoteIp := tools.GetRemoteFromContext(ctx)
	if len(remoteIp) == 0 {
		return nil, errors.New("Failed find your ip")
	}

	vpnName := peerSession.VpnId
	vpn, ok := s.vpnManager.GetLocked(vpnName)
	if !ok {
		return nil, errors.New("VPN group not exists: " + vpnName.Serialize())
	}
	mtu := vpn.GetMTU(request.GetNetwork().GetMTU())
	vpn.Release()
	fmt.Printf("   * VPN: %v\n", vpnName)

	externalIps := request.GetNetwork().GetExternalIp()
	if len(externalIps) == 0 {
		if request.GetNetwork().GetExternalEnabled() {
			fmt.Println("   * try find external ips...")
			externalIps = append(externalIps, remoteIp)
		}
	}
	fmt.Printf("   * external ips: %v\n", externalIps)

	peerSession.MTU = mtu
	peerSession.HostsLine = createHostsLine(vpn.GetHostDomain(), peerSession.Hostname, request.GetServices(), peerSession.Title)
	peerSession.ExternalIp = externalIps
	peerSession.ExternalPort = port(request.GetNetwork().GetExternalPort())
	peerSession.InternalIp = request.GetNetwork().GetInternalIp()
	peerSession.InternalPort = port(request.GetNetwork().GetInternalPort())

	return &protocol.ClientInfoResponse{}, nil
}

func port(n uint32) uint32 {
	if n == 0 {
		return constants.DEFAULT_PORT_NUMBER
	} else {
		return n
	}
}

func createHostsLine(vpnname string, host string, services []string, title string) string {
	line := host + "." + vpnname + " "
	if host != strings.ToLower(host) {
		line += strings.ToLower(host) + "." + vpnname + " "
	}
	for _, s := range services {
		line += s + "." + vpnname + " "
		if s != strings.ToLower(s) {
			line += strings.ToLower(s) + "." + vpnname + " "
		}
	}
	line += "## " + title
	return line
}
