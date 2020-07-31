// +build windows

package wireguardControl

import "github.com/gongt/wireguard-config-distribute/internal/tools"

func update(ifName string, configPath string) error {
	tools.Error("update interface %s from %s", ifName, configPath)
	return nil
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
