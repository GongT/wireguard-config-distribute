#!/usr/bin/env bash

set -Eeuo pipefail

declare -r MYDIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

cd "$MYDIR/../services"
"$MYDIR/helpers/script-sender.sh" \
	"$MYDIR/helpers/install-script-systemd.sh" \
	auto-update.sh \
	systemd \
	| bash /dev/stdin normal
