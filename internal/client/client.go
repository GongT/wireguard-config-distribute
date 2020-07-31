package client

import (
	server "github.com/gongt/wireguard-config-distribute/internal/client/client.server"
	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type wgVpnStatus struct {
	controller *wireguardControl.PeersCache

	requestedAddress    string
	givenAddress        string
	interfacePrivateKey string
}

type ClientStateHolder struct {
	quitChan  chan bool
	isQuit    bool
	isRunning bool

	sessionId types.SidType
	machineId string
	server    server.ServerStatus
	vpn       wgVpnStatus

	configData oneTimeConfig
	statusData editableConfig
}

type vpnOptions interface {
	GetPerferIp() string
}

type connectionOptions interface {
	GetServer() string

	GetGrpcInsecure() bool
	GetGrpcHostname() string
	GetGrpcServerKey() string
}

func NewClient(options connectionOptions) *ClientStateHolder {
	self := ClientStateHolder{}

	self.server = server.NewGrpcClient(options.GetServer(), server.TLSOptions{
		Insecure:  options.GetGrpcInsecure(),
		Hostname:  options.GetGrpcHostname(),
		ServerKey: options.GetGrpcServerKey(),
	})

	self.quitChan = make(chan bool, 1)
	self.isQuit = false

	return &self
}

func (self *ClientStateHolder) ConfigureVPN(options vpnOptions) {
	self.vpn.requestedAddress = options.GetPerferIp()
}
func (self *ClientStateHolder) ConfigureInterface(options wireguardControl.InterfaceOptions) {
	self.vpn.controller = wireguardControl.NewPeersCache(options)
}

type configureOptions interface {
	GetMachineID() string
	GetJoinGroup() string
	GetNetworkName() string
	GetTitle() string
	GetPerferIp() string
	GetHostname() string
	GetPublicIp() string
	GetPublicIp6() string
	GetInternalIp() string
	GetListenPort() uint16
	GetIpv6Only() bool
}

func (stat *ClientStateHolder) Configure(options configureOptions) {
	stat.configData.configure(options)

	stat.machineId = options.GetMachineID()
}

func (s *ClientStateHolder) Quit() {
	if s.isQuit {
		tools.Error("Duplicate call to Client.quit()")
		return
	}
	s.isQuit = true

	s.server.Disconnect(s.isRunning, s.sessionId)

	s.quitChan <- true
}
