package wireguardControl

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl/child_process"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type nativeInterface struct {
	lastip     string
	lastKey    string
	lastSubnet uint16
	lastPort   uint16
	lastMTU    uint16
}

var updateCmd string = "syncconf"

func generateFiltered(file string) string {
	filtered := file + ".f"
	output := child_process.RunGetStandardOutput("strip config file", "wg-quick", "strip", file)
	err := ioutil.WriteFile(filtered, []byte(output), os.FileMode(0600))
	if err != nil {
		tools.Die("Failed write filtered config file: %s", err.Error())
	}
	return filtered
}

func (wc *WireguardControl) updateInterface() error {
	tools.Debug("update interface %s from %s", wc.interfaceName, wc.configFile)
	if wc.nativeInterface == nil {
		output := child_process.RunGetOutput("check interface exists", "wg", "show", wc.interfaceName)
		if !strings.Contains(output, "No such device") {
			tools.Error("interface already exists, deleting...")
			child_process.MustSuccess("delete interface", "ip", "link", "del", "dev", wc.interfaceName)
		}

		wc.nativeInterface = &nativeInterface{}
		wc._flush()
		child_process.MustSuccess("init interface config", "wg-quick", "up", wc.configFile)
	} else {
		if wc.givenAddress == wc.nativeInterface.lastip &&
			wc.privateKey == wc.nativeInterface.lastKey &&
			wc.subnet == wc.nativeInterface.lastSubnet &&
			wc.interfaceListenPort == wc.nativeInterface.lastPort &&
			wc.interfaceMTU == wc.nativeInterface.lastMTU {
			configFileFiltered := generateFiltered(wc.configFile)
			child_process.MustSuccess("update interface config", "wg", updateCmd, wc.interfaceName, configFileFiltered)
		} else {
			child_process.MustSuccess("init interface config", "wg-quick", "down", wc.configFile)
			child_process.MustSuccess("init interface config", "wg-quick", "up", wc.configFile)
			wc._flush()
		}
	}

	return nil
}

func (wc *WireguardControl) _flush() {
	wc.nativeInterface.lastip = wc.givenAddress
	wc.nativeInterface.lastKey = wc.privateKey
	wc.nativeInterface.lastSubnet = wc.subnet
	wc.nativeInterface.lastPort = wc.interfaceListenPort
	wc.nativeInterface.lastMTU = wc.interfaceMTU
}

func (wc *WireguardControl) deleteInterface() error {
	if wc.nativeInterface != nil {
		child_process.ShouldSuccess("delete interface", "wg-quick", "down", wc.configFile)
		wc.nativeInterface = nil
	}
	return nil
}

func init() {
	tools.EnsureCommandExists("wg-quick")
	tools.EnsureCommandExists("wg")
	tools.EnsureCommandExists("ip")

	output := child_process.RunGetOutput("detect wg version", "wg", "syncconf", "-h")
	if !strings.Contains(output, "syncconf") {
		tools.Error("Your wireguard is old, please consider update.")
		updateCmd = "setconf"
	}
}
