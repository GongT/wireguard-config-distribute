#!/usr/bin/env bash

set -Eeuo pipefail

export RHOST="$1"
shift

declare -r MYDIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

cd "$MYDIR/../services"
"$MYDIR/helpers/script-sender.sh" \
	"$MYDIR/helpers/install-script-server.sh" \
	auto-update.sh \
	systemd \
	| ssh "$RHOST" bash /dev/stdin "$@"
