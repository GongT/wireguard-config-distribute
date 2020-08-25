package peerStatus

import (
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type PeerData struct {
	sessionId   types.SidType
	MachineId   string
	VpnId       types.VpnIdType
	Title       string
	Hostname    string
	PublicKey   string
	VpnIp       string
	WorkgroupId string

	MTU          uint32
	HostsLine    string
	ExternalIp   []string
	ExternalPort uint32
	InternalIp   string
	InternalPort uint32

	sender        *protocol.WireguardApi_StartServer
	lastKeepAlive time.Time
}

func (peer *PeerData) CreateId() string {
	return peer.VpnId.Serialize() + "::" + peer.WorkgroupId + "::" + peer.MachineId
}
