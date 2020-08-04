package peerStatus

import (
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/asyncChannels"
	"github.com/gongt/wireguard-config-distribute/internal/debugLocker"
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type lPeerData *PeerData

type PeerData struct {
	sessionId types.SidType
	MachineId string
	Title     string
	Hostname  string
	PublicKey string
	VpnIp     string
	MTU       uint32
	Hosts     []string

	NetworkId    string
	ExternalIp   []string
	ExternalPort uint32
	InternalIp   string
	InternalPort uint32

	sender        *protocol.WireguardApi_StartServer
	lastKeepAlive time.Time
}

type PeerStatus struct {
	list     map[types.SidType]lPeerData
	m        debugLocker.MyLocker
	onChange *asyncChannels.AsyncChanSidType

	guid    types.SidType
	guidMap map[string]types.SidType
}

func NewPeerStatus() *PeerStatus {
	return &PeerStatus{
		list:     make(map[types.SidType]lPeerData),
		onChange: asyncChannels.NewChanSidType(),
		m:        debugLocker.NewMutex(),

		guid:    0,
		guidMap: make(map[string]types.SidType),
	}
}
