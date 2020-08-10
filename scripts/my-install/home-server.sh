#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"
export RHOST="home.gongt.me"

pwsh scripts/build.ps1 client

echo
echo

function x() {
	echo "$*"
	"$@"
}

x rsync --progress \
	scripts/services/client@.service \
	dist/client \
	$RHOST:/tmp/
	
cat <<- 'EOF' | ssh $RHOST bash
	set -xEeuo pipefail

	systemctl stop wireguard-config-client@normal || true
	rm -f /usr/local/bin/wireguard-config-client

	cp /tmp/client@.service /usr/lib/systemd/system/wireguard-config-client@.service
	systemctl daemon-reload
	systemctl enable wireguard-config-client@normal

	cp /tmp/client /usr/local/bin/wireguard-config-client
	systemctl restart wireguard-config-client@normal
