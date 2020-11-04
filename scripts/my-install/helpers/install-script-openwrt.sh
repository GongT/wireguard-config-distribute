#!/usr/bin/env bash

set -Eeuo pipefail

function die() {
	echo "$*" >&2
	exit 1
}

if [[ "$(basename "$(readlink -e /proc/1/exe)")" != "procd" ]]; then
	die "remote server did not running procd, please check environment."
fi

mkdir -p /usr/local/libexec/wireguard-config-client
cp auto-update.sh openwrt/service-control.sh /usr/local/libexec/wireguard-config-client
cp openwrt/procd-init.sh /etc/init.d/wireguard-config-client

UPDATE_SCRIPT="/usr/local/libexec/wireguard-config-client/auto-update.sh"
chmod a+x "$UPDATE_SCRIPT"

{
	crontab -l | grep -v ' wireguard-config-client-auto-update '
	echo "0 0 * * * bash '$UPDATE_SCRIPT' client #! wireguard-config-client-auto-update !#"
} | crontab -

if ! [[ -e /usr/local/libexec/wireguard-config-client/client.alpine ]]; then
	bash "$UPDATE_SCRIPT" client
else
	/etc/init.d/wireguard-config-client restart
fi
