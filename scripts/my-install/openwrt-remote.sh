#!/usr/bin/env bash

set -Eeuo pipefail

export RHOST="$1"

declare -r MYDIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

cd "$MYDIR/../services"
"$MYDIR/helpers/script-sender.sh" \
	"$MYDIR/helpers/install-script-openwrt.sh" \
	auto-update.sh \
	openwrt \
	| ssh "$RHOST" bash
	