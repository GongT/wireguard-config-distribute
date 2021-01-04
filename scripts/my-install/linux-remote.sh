#!/usr/bin/env bash

set -Eeuo pipefail

export RHOST="$1"

declare -r MYDIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

cd "$MYDIR/../services"
"$MYDIR/helpers/script-sender.sh" \
	"$MYDIR/helpers/install-script-systemd.sh" \
	auto-update.sh \
	systemd \
	| ssh -T -o "SendEnv" "$RHOST" bash "<(echo $(base64 -w0 "$MYDIR/helpers/script-receiver.sh") | base64 -d)" normal
