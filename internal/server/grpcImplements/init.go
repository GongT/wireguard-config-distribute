package grpcImplements

import (
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements/networkManager"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements/peerStatus"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements/vpnManager"
	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
)

type ServerImplementOptions interface {
	GetPassword() string
	GetStorageLocation() string
	GetGrpcInsecure() bool
}

type serverImplement struct {
	password string
	storage  *storage.ServerStorage
	insecure bool

	peerStatus     *peerStatus.PeerStatus
	networkManager *networkManager.NetworkManager
	vpnManager     *vpnManager.VpnManager
}

func CreateServerImplement(opts ServerImplementOptions) *serverImplement {
	store := storage.CreateStorage(opts.GetStorageLocation())
	return &serverImplement{
		password:       opts.GetPassword(),
		storage:        store,
		insecure:       opts.GetGrpcInsecure(),
		peerStatus:     peerStatus.NewPeerStatus(),
		networkManager: networkManager.NewNetworkManager(store),
		vpnManager:     vpnManager.NewVpnManager(store),
	}
}
