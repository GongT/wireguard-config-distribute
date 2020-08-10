package wgexe

import (
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl/child_process"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type wgExe struct {
	updateCmd string
}

var wg_cli_cache *wgExe

func GetWireguardCli() *wgExe {
	if wg_cli_cache != nil {
		return wg_cli_cache
	}

	tools.EnsureCommandExists(WG_BINARY, INSTALL_INFORMATION)

	wg_cli_cache = &wgExe{}

	output := child_process.RunGetOutput("detect wg version", WG_BINARY, "syncconf", "-h")
	if strings.Contains(output, "Usage: wg setconf") {
		tools.Error(`Good! Your wireguard is not very old.`)
		wg_cli_cache.updateCmd = "syncconf"
	} else {
		tools.Error(`===========================================
%s
Your wireguard is old, please consider update.
===========================================`, output)
		wg_cli_cache.updateCmd = "setconf"
	}
	return wg_cli_cache
}

func (wg *wgExe) SmallChange(interfaceName, configFileFiltered string) error {
	return child_process.ShouldSuccess("update interface config", WG_BINARY, wg.updateCmd, interfaceName, configFileFiltered)
}

func (wg *wgExe) Exists(interfaceName string) bool {
	output := child_process.RunGetOutput("check interface exists", WG_BINARY, "show", interfaceName)
	return !strings.Contains(output, NO_DEVICE_STRING)
}
