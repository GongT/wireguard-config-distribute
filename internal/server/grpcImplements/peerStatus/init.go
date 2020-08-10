package peerStatus

import (
	"github.com/gongt/wireguard-config-distribute/internal/asyncChannels"
	"github.com/gongt/wireguard-config-distribute/internal/debugLocker"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type vpnPeersMap map[types.SidType]*PeerData

type PeersManager struct {
	list   map[types.SidType]*PeerData
	mapper map[types.VpnIdType]vpnPeersMap

	m        debugLocker.MyLocker
	onChange *asyncChannels.AsyncChanSidType

	guid    types.SidType
	guidMap map[string]types.SidType
}

func NewPeersManager() *PeersManager {
	return &PeersManager{
		list:     make(map[types.SidType]*PeerData),
		mapper:   make(map[types.VpnIdType]vpnPeersMap),
		onChange: asyncChannels.NewChanSidType(),
		m:        debugLocker.NewMutex(),

		guid:    0,
		guidMap: make(map[string]types.SidType),
	}
}
