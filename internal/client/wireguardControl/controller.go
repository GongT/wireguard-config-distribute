package wireguardControl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/debugLocker"
)

type WireguardControl struct {
	interfaceName   string
	nativeInterface *nativeInterface

	peers      []peerData
	configFile string

	mu debugLocker.MyLocker

	requestedAddress string
	givenAddress     string
	privateKey       string
	subnet           uint16

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
}

func NewWireguardControl(options VpnOptions) *WireguardControl {
	dir := getTempDir()

	return &WireguardControl{
		interfaceName: options.GetInterfaceName(),

		peers:      make([]peerData, 20),
		configFile: filepath.Join(dir, options.GetInterfaceName()+".conf"),

		requestedAddress: options.GetPerferIp(),
		givenAddress:     "",
		privateKey:       "",

		interfaceTitle:      fmt.Sprintf("%s (%s)", options.GetHostname(), options.GetTitle()),
		interfaceListenPort: options.GetListenPort(),
		interfaceMTU:        options.GetMTU(),

		mu: debugLocker.NewMutex(),
	}
}

func (wc *WireguardControl) creatConfigFile() error {
	return ioutil.WriteFile(wc.configFile, wc.creatConfig(), os.FileMode(0600))
}

func (wc *WireguardControl) DeleteInterface() error {
	return wc.deleteInterface()
}
