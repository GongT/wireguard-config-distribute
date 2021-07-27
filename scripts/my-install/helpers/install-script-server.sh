#!/usr/bin/env bash

set -Eeuo pipefail

function die() {
	echo "$*" >&2
	exit 1
}

if [[ "$(basename "$(readlink -e /proc/1/exe)")" != "systemd" ]]; then
	die "remote server did not running systemd, please check environment."
fi

mkdir -p /usr/local/libexec/wireguard-config-client

echo "installing systemd service"
{
	sed '/\[Service]/Q' systemd/server.service
	echo '[Service]'
	echo "ExecStartPre=-+/usr/local/libexec/wireguard-config-client/auto-update.sh server"
	sed '1,/\[Service]/d' systemd/server.service
} >/usr/lib/systemd/system/wireguard-config-server.service

echo "reload"
systemctl daemon-reload

echo "install server application"
bash ./systemd/install-update-service.sh server
DISABLE_RESTART=yes bash auto-update.sh server

if [[ $* == *--restart* ]]; then
	echo "enable and restart server"
	systemctl enable wireguard-config-server.service
	systemctl restart wireguard-config-server.service
else
	echo "enable and start server (restart with --restart)"
	systemctl enable --now wireguard-config-server.service
fi
