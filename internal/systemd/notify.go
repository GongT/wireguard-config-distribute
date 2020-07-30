// +build linux

package systemd

import (
	"github.com/coreos/go-systemd/daemon"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func ChangeToReady() {
	_, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if err != nil {
		tools.Error("Failed send systemd event: %s", err.Error())
	}
}

func ChangeToQuit() {
	_, err := daemon.SdNotify(false, daemon.SdNotifyStopping)
	if err != nil {
		tools.Error("Failed send systemd event: %s", err.Error())
	}
}

func UpdateState(status string) {
	daemon.SdNotify(false, status)
}
