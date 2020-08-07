package wireguardControl

import (
	"fmt"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl/interfaceState"
	"github.com/gongt/wireguard-config-distribute/internal/debugLocker"
)

type WireguardControl struct {
	interfaceName   string
	nativeInterface interfaceState.InterfaceState
	dryRun          bool

	peers      []peerData
	configFile string

	extendedConfigCreated bool

	mu debugLocker.MyLocker

	requestedAddress string
	givenAddress     string
	privateKey       string
	subnet           uint8

	interfaceTitle      string
	interfaceListenPort uint16
	interfaceMTU        uint16
}

type VpnOptions interface {
	GetPerferIp() string

	GetListenPort() uint16
	GetInterfaceName() string
	GetMTU() uint16
	GetTitle() string
	GetHostname() string

	GetNetworkName() string

	GetDryRun() bool
}

func NewWireguardControl(options VpnOptions) *WireguardControl {
	var nativeInterface interfaceState.InterfaceState
	if options.GetDryRun() {
		nativeInterface = interfaceState.CreateDummy()
	} else {
		nativeInterface = interfaceState.CreateInterface(options.GetInterfaceName())
	}
	return &WireguardControl{
		interfaceName: options.GetInterfaceName(),

		nativeInterface: nativeInterface,
		dryRun:          options.GetDryRun(),

		peers:      make([]peerData, 20),
		configFile: filepath.Join(TempDir, options.GetInterfaceName()+".native.conf"),

		extendedConfigCreated: false,

		requestedAddress: options.GetPerferIp(),
		givenAddress:     "",
		privateKey:       "",

		interfaceTitle:      fmt.Sprintf("%s (%s) [AT] %s", options.GetHostname(), options.GetTitle(), options.GetNetworkName()),
		interfaceListenPort: options.GetListenPort(),
		interfaceMTU:        options.GetMTU(),

		mu: debugLocker.NewMutex(),
	}
}

func (wc *WireguardControl) DeleteInterface() error {
	return wc.nativeInterface.DeleteInterface()
}
