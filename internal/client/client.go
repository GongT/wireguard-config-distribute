package client

import (
	server "github.com/gongt/wireguard-config-distribute/internal/client/client.server"
	"github.com/gongt/wireguard-config-distribute/internal/client/sharedConfig"
	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type ClientStateHolder struct {
	quitChan  chan bool
	isQuit    bool
	isRunning bool

	sessionId types.SidType
	machineId string
	server    *server.ServerStatus
	vpn       *wireguardControl.WireguardControl

	configData oneTimeConfig
	statusData editableConfig

	password string
}

func NewClient(options sharedConfig.ReadOnlyConnectionOptions) *ClientStateHolder {
	self := ClientStateHolder{}

	self.server = server.NewGrpcClient(options.GetServer(), options.GetPassword(), options)

	self.quitChan = make(chan bool, 1)
	self.isQuit = false

	return &self
}

func (self *ClientStateHolder) ConfigureVPN(options wireguardControl.VpnOptions) {
	self.vpn = wireguardControl.NewWireguardControl(options)
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

	s.vpn.DeleteInterface()
	s.server.Disconnect(s.isRunning, s.sessionId)

	s.quitChan <- true
}
