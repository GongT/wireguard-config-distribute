package client

import (
	server "github.com/gongt/wireguard-config-distribute/internal/client/client.server"
	"github.com/gongt/wireguard-config-distribute/internal/client/sharedConfig"
	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/transport"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type ClientStateHolder struct {
	quitChan  chan bool
	isQuit    bool
	isRunning bool

	ipv4Only bool

	sessionId types.SidType
	machineId string
	server    *server.ServerStatus
	vpn       *wireguardControl.WireguardControl
	nat       *transport.Transport

	configData oneTimeConfig
	statusData editableConfig

	password string

	hostsHandler HandlerFunction
}

func NewClient(options sharedConfig.ReadOnlyConnectionOptions) *ClientStateHolder {
	self := ClientStateHolder{}

	self.server = server.NewGrpcClient(options.GetServer(), options.GetPassword(), options)
	self.nat = transport.NewTransport()

	self.quitChan = make(chan bool, 1)
	self.isQuit = false

	return &self
}

type configureOptions interface {
	wireguardControl.VpnOptions

	GetMachineID() string
	GetJoinGroup() string
	GetPublicIp() string
	GetPublicIp6() string
	GetInternalIp() string
	GetListenPort() uint16
	GetPublicPort() uint16
	GetIpv6Only() bool
	GetMTU() uint16

	GetIpv4Only() bool
}

func (stat *ClientStateHolder) Configure(options configureOptions) {
	stat.configData.configure(options)

	stat.vpn = wireguardControl.NewWireguardControl(options)
	stat.ipv4Only = options.GetIpv4Only()

	stat.machineId = options.GetMachineID()
}

func (s *ClientStateHolder) Quit() {
	if s.isQuit {
		tools.Error("Duplicate call to Client.quit()")
		return
	}
	s.isQuit = true

	s.nat.Quit()

	if s.vpn != nil {
		tools.Debug("deleting wg interface")
		s.vpn.DeleteInterface()
	}
	tools.Debug("wg interface down")

	tools.Error("disconnect grpc")
	s.server.Disconnect(s.isRunning, s.sessionId)
	tools.Error("grpc is end")

	s.quitChan <- true
	close(s.quitChan)
}
