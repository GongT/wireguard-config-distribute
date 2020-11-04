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
	sed '/\[Service]/Q' systemd/client@.service
	echo "Wants=wireguard-config-auto-update.service"
	echo "After=wireguard-config-auto-update.service"
	echo '[Service]'
	sed '1,/\[Service]/d' systemd/client@.service
} >/usr/lib/systemd/system/wireguard-config-client@.service

touch "/usr/local/libexec/wireguard-config-client/client-should-update"
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

rm -f /usr/local/libexec/ensure-kmod.sh

cp systemd/service-control.sh systemd/ensure-kmod.sh auto-update.sh /usr/local/libexec/wireguard-config-client
chmod a+x /usr/local/libexec/wireguard-config-client/auto-update.sh

START=$(date +%s)

systemctl daemon-reload

SERVICES=()
for I; do
	SERVICES+=("wireguard-config-client@$I.service")
done

systemctl enable wireguard-config-auto-update.timer
systemctl restart wireguard-config-auto-update.timer

systemctl stop wireguard-config-auto-update.service
systemctl enable --now "${SERVICES[@]}"

function checkStatus() {
	echo "==========="
	SERVICES+=(wireguard-config-auto-update.service)

	local -a FAILEDLIST=() ARGS=()
	local I IID
	for I in "${SERVICES[@]}"; do
		if ! systemctl is-active --quiet "$I"; then
			FAILEDLIST+=("$I")
			ARGS+=("--unit=$I")
		fi
	done

	local -i delta=$(($(date +%s) - START))
	# delta=$((delta + 1))

	echo journalctl --no-legend --no-pager "${ARGS[@]}" --since "$delta second ago"
	journalctl --no-pager "${ARGS[@]}" --since "$delta second ago"
	echo "service failed start" >&2
	for I in "${FAILEDLIST[@]}"; do
		echo -e "  *\e[38;5;9m $I\e[0m"
	done
}
