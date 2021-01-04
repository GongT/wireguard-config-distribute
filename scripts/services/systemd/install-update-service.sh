#!/usr/bin/env bash

set -Eeuo pipefail

declare -r TYPE=$1
declare -r distRoot="/usr/local/libexec/wireguard-config-client"

touch "$distRoot/$TYPE-should-update"
{
	sed '/\[Service]/Q' systemd/auto-update.service
	echo '[Service]'
	if [[ -e "$distRoot/client-should-update" ]]; then
		echo "ExecStart=$distRoot/auto-update.sh client"
	fi
	if [[ -e "$distRoot/server-should-update" ]]; then
		echo "ExecStart=$distRoot/auto-update.sh server"
	fi
	sed '1,/\[Service]/d' systemd/auto-update.service
} >/usr/lib/systemd/system/wireguard-config-auto-update.service

cp systemd/auto-update.timer /usr/lib/systemd/system/wireguard-config-auto-update.timer

cp auto-update.sh $distRoot
chmod a+x $distRoot/auto-update.sh

systemctl daemon-reload
systemctl stop wireguard-config-auto-update.service &>/dev/null

if [[ -e "$distRoot/client-should-update" ]] && ! [[ -e "$distRoot/client" ]]; then
	$distRoot/auto-update.sh client
fi
if [[ -e "$distRoot/server-should-update" ]] && ! [[ -e "$distRoot/server" ]]; then
	$distRoot/auto-update.sh server
fi

systemctl enable wireguard-config-auto-update.timer
systemctl restart wireguard-config-auto-update.timer
