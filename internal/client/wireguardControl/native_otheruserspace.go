// +build aix android darwin dragonfly freebsd illumos js netbsd openbsd plan9 solaris

package wireguardControl

import "github.com/gongt/wireguard-config-distribute/internal/tools"

func (wc *WireguardControl) updateInterface() error {
	tools.Error("update interface %s from %s", wc.interfaceName, wc.configFile)
	return nil
}

func (wc *WireguardControl) deleteInterface() error {
	return nil
}
