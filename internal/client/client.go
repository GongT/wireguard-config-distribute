package client

import (
	server "github.com/gongt/wireguard-config-distribute/internal/client/client.server"
	"github.com/gongt/wireguard-config-distribute/internal/client/sharedConfig"
	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl"
	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/transport"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

type ClientStateHolder struct {
	quitChan  chan bool
	isQuit    bool
	qDispose  func()
	isRunning bool

	ipDetect *detect_ip.Detect

	sessionId types.SidType
	machineId string
	server    *server.ServerStatus
	vpn       *wireguardControl.WireguardControl
	nat       *transport.Transport

	privateStatus infoConfig
	sharedStatus  editableConfig

	password string

	hostsHandler HandlerFunction
}

func NewClient(options sharedConfig.ReadOnlyConnectionOptions) *ClientStateHolder {
	self := ClientStateHolder{}

	self.server = server.NewGrpcClient(options.GetServer(), options.GetPassword(), options)
	self.nat = transport.NewTransport()

	self.quitChan = make(chan bool, 1)
	self.isQuit = false
	self.qDispose = tools.WaitExit(func(int) {
		self.Quit()
	})

	return &self
}

type configureOptions interface {
	wireguardControl.VpnOptions

	infoOptions

	GetMachineID() string

	detect_ip.Options
}

func (stat *ClientStateHolder) Configure(options configureOptions) {
	stat.privateStatus.configure(options)

	stat.vpn = wireguardControl.NewWireguardControl(options, stat.privateStatus.createInterfaceComment())

	stat.ipDetect = detect_ip.NewDetect(options)

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
