#!/usr/bin/env bash

set -Eeuo pipefail

function die() {
	echo "$*" >&2
	exit 1
}

if [[ "$(basename "$(readlink -e /proc/1/exe)")" != "systemd" ]]; then
	die "remote server did not running systemd, please check environment."
fi

rm -rf /usr/local/libexec/wireguard-config-client

mkdir -p /usr/local/libexec/wireguard-config
echo "copy server files" >&2
cp systemd/service-control.sh systemd/ensure-kmod.sh /usr/local/libexec/wireguard-config

echo "installing systemd service"
cp systemd/server.service /usr/lib/systemd/system/wireguard-config-server.service

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
