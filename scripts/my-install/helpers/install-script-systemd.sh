#!/usr/bin/env bash

set -Eeuo pipefail

function die() {
	echo "$*" >&2
	exit 1
}

if [[ "$(basename "$(readlink -e /proc/1/exe)")" != "systemd" ]]; then
	die "remote server did not running systemd, please check environment."
fi
if ! command wg &>/dev/null; then
	die "wireguard tools (wg binary) is not found"
fi

mkdir -p /usr/local/libexec/wireguard-config-client

echo "create wireguard-config-client@.service" >&2
{
	sed '/\[Service]/Q' systemd/client@.service
	echo '[Service]'
	echo "ExecStartPre=-+/usr/local/libexec/wireguard-config-client/auto-update.sh client"
	sed '1,/\[Service]/d' systemd/client@.service
} >/usr/lib/systemd/system/wireguard-config-client@.service

echo "install update service" >&2
bash ./systemd/install-update-service.sh client

rm -f /usr/local/libexec/ensure-kmod.sh

echo "copy client files" >&2
cp systemd/service-control.sh systemd/ensure-kmod.sh /usr/local/libexec/wireguard-config-client

START=$(date +%s)

systemctl daemon-reload

SERVICES=()
for I; do
	SERVICES+=("wireguard-config-client@$I.service")
done

echo "enable services: ${SERVICES[*]}" >&2
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
