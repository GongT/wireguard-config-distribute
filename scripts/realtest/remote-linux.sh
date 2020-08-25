#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

echo -ne "\ec"

export GOOS="linux"
export GOARCH="amd64"
export RHOST="$1"

pwsh scripts/build.ps1 client

rsync -v dist/client $RHOST:/tmp/wireguard-config-client

ssh -tt $RHOST bash -c "
	set -a
	source /etc/wireguard/client.normal.conf
	source /etc/wireguard/client.conf
	set +a
	exec /tmp/wireguard-config-client --group normal
"
