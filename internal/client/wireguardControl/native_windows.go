package wireguardControl

import (
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type nativeInterface struct {
}

func (wc *WireguardControl) updateInterface() error {
	tools.Error("update interface %s from %s", wc.interfaceName, wc.configFile)
	return nil
}

func (wc *WireguardControl) deleteInterface() error {
	return nil
}

func init() {
	tools.EnsureCommandExists(`C:\Program Files\WireGuard\wireguard.exe`)
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
