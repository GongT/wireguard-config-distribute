// +build linux

package systemd

import (
	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func ChangeToReady() {
	tools.Error("[systemd] READY")
	_, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if err != nil {
		tools.Error("Failed send systemd event: %s", err.Error())
	}
}

func ChangeToReload() {
	tools.Error("[systemd] RELOAD")
	_, err := daemon.SdNotify(false, daemon.SdNotifyReloading)
	if err != nil {
		tools.Error("Failed send systemd event: %s", err.Error())
	}
}

func ChangeToQuit() {
	tools.Error("[systemd] STOP")
	_, err := daemon.SdNotify(false, daemon.SdNotifyStopping)
	if err != nil {
		tools.Error("Failed send systemd event: %s", err.Error())
	}
}

func UpdateState(status string) {
	tools.Error("[systemd] STATE: " + status)
	daemon.SdNotify(false, status)
}
