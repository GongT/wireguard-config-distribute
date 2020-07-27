package client

import (
	server "github.com/gongt/wireguard-config-distribute/internal/client/client.server"
)

type wgVpnStatus struct {
	requestedAddress    string
	givenAddress        string
	interfacePrivateKey string
}

type clientStateHolder struct {
	quitChan  chan bool
	isQuit    bool
	isRunning bool

	SessionId uint64
	server    server.ServerStatus
	vpn       wgVpnStatus

	configData oneTimeConfig
	statusData editableConfig
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

type configureOptions interface {
	GetJoinGroup() string
	GetNetworkName() string
	GetTitle() string
	GetPerferIp() string
	GetHostname() string
	GetPublicIp() string
	GetPublicIp6() string
	GetInternalIp() []string
	GetListenPort() uint16
	GetIpv6Only() bool
}

func (stat *clientStateHolder) Configure(options configureOptions) {
	stat.configData.configure(options)
}
