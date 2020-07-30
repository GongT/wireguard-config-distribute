// +build !windows

package wireguardControl

import "github.com/gongt/wireguard-config-distribute/internal/tools"

func update(ifName string, configPath string) error {
	tools.Error("update interface %s from %s", ifName, configPath)
	return nil
}
