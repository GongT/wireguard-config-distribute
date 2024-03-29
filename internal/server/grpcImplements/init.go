package grpcImplements

import (
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements/peerStatus"
	"github.com/gongt/wireguard-config-distribute/internal/server/grpcImplements/vpnManager"
	"github.com/gongt/wireguard-config-distribute/internal/server/storage"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type ServerImplementOptions interface {
	GetStorageLocation() string
	GetGrpcInsecure() bool
}

type PeerObject struct {
	Data peerStatus.PeerData
}

type Implements struct {
	protocol.UnimplementedWireguardApiServer

	storage  *storage.ServerStorage
	insecure bool
	isQuit   bool

	vpnManager   *vpnManager.VpnManager
	peersManager *peerStatus.PeersManager

	keepAliveTimer *time.Ticker
	quitCh         chan bool
}

func CreateServerImplement(opts ServerImplementOptions) *Implements {
	store := storage.CreateStorage(opts.GetStorageLocation())

	srv := Implements{
		storage:  store,
		insecure: opts.GetGrpcInsecure(),
		isQuit:   false,

		vpnManager:   vpnManager.NewVpnManager(store),
		peersManager: peerStatus.NewPeersManager(),

		keepAliveTimer: nil,
		quitCh:         make(chan bool, 1),
	}

	return &srv
}

func (srv *Implements) StartWorker() {
	srv.keepAliveTimer = time.NewTicker(1 * time.Minute)
	go srv.peersManager.StartHandleChange()
	go func() {
		for {
			select {
			case <-srv.keepAliveTimer.C:
				srv.peersManager.CleanupTimeoutPeers()
			case <-srv.quitCh:
				return
			}
		}
	}()
}

func (srv *Implements) Quit() {
	if srv.isQuit {
		tools.Error("Duplicate call to Implements.Stop()")
		return
	}
	srv.isQuit = true

	srv.peersManager.StopHandleChange()
	srv.keepAliveTimer.Stop()
	srv.quitCh <- true
	close(srv.quitCh)
}
