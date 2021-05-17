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

{
	sed '/\[Service]/Q' systemd/server.service
	echo '[Service]'
	echo "ExecStartPre=-+/usr/local/libexec/wireguard-config-client/auto-update.sh server"
	sed '1,/\[Service]/d' systemd/server.service
} >/usr/lib/systemd/system/wireguard-config-server.service

bash ./systemd/install-update-service.sh server

systemctl enable --now wireguard-config-server.service
