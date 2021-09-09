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

rm -rf /usr/local/libexec/wireguard-config-client /etc/init.d/wireguard-config-*

mkdir -p /usr/local/libexec/wireguard-config
cp auto-update.sh openwrt/service-control.sh /usr/local/libexec/wireguard-config
cp openwrt/procd-init.sh /etc/init.d/wireguard-config

UPDATE_SCRIPT="/usr/local/libexec/wireguard-config/auto-update.sh"
chmod a+x "$UPDATE_SCRIPT"

{
	crontab -l | grep -v ' wireguard-config-client-auto-update '
	echo "0 0 * * * bash '$UPDATE_SCRIPT' client #! wireguard-config-client-auto-update !#"
} | crontab -

if ! [[ -e /usr/local/libexec/wireguard-config/client.alpine ]]; then
	bash "$UPDATE_SCRIPT" client
else
	/etc/init.d/wireguard-config restart
fi
