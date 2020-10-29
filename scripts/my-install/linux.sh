#!/usr/bin/env bash

set -Eeuo pipefail

declare -r MYDIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

bash "$MYDIR/linux/download-script.sh"

cd "$MYDIR/../services"
"$MYDIR/linux/script-sender.sh" \
	"$MYDIR/linux/install-script-systemd.sh" \
	client@.service \
	auto-update.service \
	auto-update.timer \
	auto-update.sh \
	ensure-kmod.sh \
	| bash -c "$(<"$MYDIR/linux/script-receiver.sh")" -- normal
