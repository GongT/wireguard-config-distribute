#!/usr/bin/env bash

set -Eeuo pipefail

RHOST=gateway.gongt.me

declare -r MYDIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd "$MYDIR/.."

pwsh build.ps1 server

cd ".."

"$MYDIR/helpers/script-sender.sh" \
	"$MYDIR/helpers/test-script-server.sh" \
	dist/server \
	| ssh "$RHOST" bash /dev/stdin
