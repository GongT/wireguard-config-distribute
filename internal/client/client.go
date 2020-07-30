package client

import (
	server "github.com/gongt/wireguard-config-distribute/internal/client/client.server"
	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type wgVpnStatus struct {
	controller *wireguardControl.PeersCache

	requestedAddress    string
	givenAddress        string
	interfacePrivateKey string
}

type clientStateHolder struct {
	quitChan  chan bool
	isQuit    bool
	isRunning bool

	MachineId string
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

func NewClient(options connectionOptions) *clientStateHolder {
	self := clientStateHolder{}

	self.server = server.NewGrpcClient(options.GetServer(), server.TLSOptions{
		Insecure:  options.GetGrpcInsecure(),
		Hostname:  options.GetGrpcHostname(),
		ServerKey: options.GetGrpcServerKey(),
	})

	self.quitChan = make(chan bool, 1)
	self.isQuit = false

	return &self
}

func (self *clientStateHolder) ConfigureVPN(options vpnOptions) {
	self.vpn.requestedAddress = options.GetPerferIp()
}
func (self *clientStateHolder) ConfigureInterface(options wireguardControl.InterfaceOptions) {
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

func (stat *clientStateHolder) Configure(options configureOptions) {
	stat.configData.configure(options)

	stat.MachineId = options.GetMachineID()
}

func (s *clientStateHolder) Quit() {
	if s.isQuit {
		tools.Error("Duplicate call to Client.quit()")
		return
	}
	s.isQuit = true

	s.server.Disconnect(s.isRunning, s.MachineId)

	s.quitChan <- true
}
