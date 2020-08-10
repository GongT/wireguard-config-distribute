package interfaceState

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl/child_process"
	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl/wgexe"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type InterfaceOptions interface {
	GetAddress() string
	GetNetwork() string
	GetMtu() int
	GetConfigFile() (string, error)
}

const WG_MANAGE = `C:\Program Files\WireGuard\wireguard.exe`

type nativeState struct {
	tunnelServiceName string
}

func (is *interfaceState) init() {
	is.native = &nativeState{
		tunnelServiceName: "WireGuardTunnel$" + is.ifname,
	}
	tools.EnsureCommandExists(WG_MANAGE, "You should install wireguard from https://download.wireguard.com/windows-client/")
}

func (is *interfaceState) DeleteInterface() error {
	if checkServiceExists(is.native.tunnelServiceName) {
		err := child_process.ShouldSuccess("delete wg service", WG_MANAGE, "/uninstalltunnelservice", is.ifname)
		if err != nil {
			return fmt.Errorf("failed delete: %v", err)
		}

		is.network = ""
		is.mtu = 0

		time.Sleep(2 * time.Second)
		for wgexe.GetWireguardCli().Exists(is.ifname) {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (is *interfaceState) CreateOrUpdateInterface(options InterfaceOptions) error {
	changed := diffState(is, options)
	if changed.network || changed.mtu {
		if err := is.DeleteInterface(); err != nil {
			return fmt.Errorf("failed create or update: %v", err)
		}

		path, err := options.GetConfigFile()
		if err != nil {
			return fmt.Errorf("failed create or update: %v", err)
		}

		child_process.MustSuccess("install wg service", WG_MANAGE, "/installtunnelservice", path)

		time.Sleep(2 * time.Second)
		for !wgexe.GetWireguardCli().Exists(is.ifname) {
			time.Sleep(1 * time.Second)
		}
	}

	changed.commit()
	return nil
}

func checkServiceExists(name string) bool {
	return child_process.RunGetReturnCode("check service exists", "sc.exe", "query", name) == 0
}

/*
用法: C:\Program Files\WireGuard\wireguard.exe [
    (无参数)：提升并安装管理服务
    /installmanagerservice
    /installtunnelservice CONFIG_PATH
    /uninstallmanagerservice
    /uninstalltunnelservice TUNNEL_NAME
    /managerservice
    /tunnelservice CONFIG_PATH
    /ui CMD_READ_HANDLE CMD_WRITE_HANDLE CMD_EVENT_HANDLE LOG_MAPPING_HANDLE
    /dumplog OUTPUT_PATH
    /update [LOG_FILE]
]
*/
