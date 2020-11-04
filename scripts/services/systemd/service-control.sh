#!/usr/bin/env bash

set -Eeuo pipefail

declare -r ACTION=$1

function die() {
	echo -e "\e[38;5;9m$*\3[0m"
	exit 1
}
function info() {
	echo -e "$*" >&2
}
function mute() {
	echo -e "\e[2m$*\e[0m" >&2
}

mapfile -t SERVICES < <(systemctl list-units --no-pager --no-legend wireguard-config-* | cut -d ' ' -f 1)
info "$ACTION services: (no block)"
for I in "${SERVICES[@]}"; do
	if [[ $I == "wireguard-config-auto-update"* ]]; then
		mute "  * $I ... (skip)"
		continue
	fi
	info "  * $I ..."
	systemctl --no-block "$ACTION" "$I"
done
