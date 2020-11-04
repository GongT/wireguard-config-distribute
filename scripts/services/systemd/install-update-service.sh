#!/usr/bin/env bash

set -Eeuo pipefail

declare -r TYPE=$1

touch "/usr/local/libexec/wireguard-config-client/$TYPE-should-update"
{
	sed '/\[Service]/Q' systemd/auto-update.service
	echo '[Service]'
	if [[ -e "/usr/local/libexec/wireguard-config-client/client-should-update" ]]; then
		echo 'ExecStart=/usr/local/libexec/wireguard-config-client/auto-update.sh client'
	fi
	if [[ -e "/usr/local/libexec/wireguard-config-client/server-should-update" ]]; then
		echo 'ExecStart=/usr/local/libexec/wireguard-config-client/auto-update.sh server'
	fi
	sed '1,/\[Service]/d' systemd/auto-update.service
} >/usr/lib/systemd/system/wireguard-config-auto-update.service

cp systemd/auto-update.timer /usr/lib/systemd/system/wireguard-config-auto-update.timer

cp auto-update.sh /usr/local/libexec/wireguard-config-client
chmod a+x /usr/local/libexec/wireguard-config-client/auto-update.sh

systemctl daemon-reload
systemctl enable wireguard-config-auto-update.timer
systemctl restart wireguard-config-auto-update.timer
systemctl stop wireguard-config-auto-update.service &>/dev/null
