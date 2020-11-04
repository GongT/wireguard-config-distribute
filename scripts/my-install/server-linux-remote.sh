#!/usr/bin/env bash

set -Eeuo pipefail

export RHOST="$1"

declare -r MYDIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

cd "$MYDIR/../services"
"$MYDIR/helpers/script-sender.sh" \
	"$MYDIR/helpers/install-script-server.sh" \
	auto-update.sh \
	systemd \
	| ssh "$RHOST" bash -c "$(<"$MYDIR/helpers/script-receiver.sh")"
