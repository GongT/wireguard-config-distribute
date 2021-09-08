#!/usr/bin/env bash

set -Eeuo pipefail

function die() {
	echo "$*" >&2
	exit 1
}

if [[ "$(basename "$(readlink -e /proc/1/exe)")" != "procd" ]]; then
	die "remote server did not running procd, please check environment."
fi
if ! command wg &>/dev/null; then
	die "wireguard tools (wg binary) is not found"
fi

mkdir -p /usr/local/libexec/wireguard-config-client
cp auto-update.sh openwrt/service-control.sh /usr/local/libexec/wireguard-config-client
cp openwrt/procd-init-server.sh /etc/init.d/wireguard-config-server

UPDATE_SCRIPT="/usr/local/libexec/wireguard-config-client/auto-update.sh"
chmod a+x "$UPDATE_SCRIPT"

{
	crontab -l | grep -v ' wireguard-config-server-auto-update '
	echo "0 0 * * * bash '$UPDATE_SCRIPT' server #! wireguard-config-server-auto-update !#"
} | crontab -

if ! [[ -e /usr/local/libexec/wireguard-config-client/server.alpine ]]; then
	bash "$UPDATE_SCRIPT" server
else
	/etc/init.d/wireguard-config-server restart
fi
