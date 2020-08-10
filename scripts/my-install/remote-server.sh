#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"
export RHOST="services.gongt.me"

pwsh scripts/build.ps1 server
pwsh scripts/build.ps1 client

echo
echo


function x() {
	echo "$*"
	"$@"
}

x rsync --progress \
	scripts/services/server.service \
	scripts/services/client@.service \
	dist/server \
	dist/client \
	$RHOST:/data/temp-images/

cat <<- 'EOF' | ssh $RHOST bash 
	set -xEeuo pipefail

	systemctl stop wireguard-config-server wireguard-config-client@service wireguard-config-client@normal || true

	rm -f /usr/local/bin/wireguard-config-server /usr/local/bin/wireguard-config-client

	cp /data/temp-images/server /usr/local/bin/wireguard-config-server
	cp /data/temp-images/client /usr/local/bin/wireguard-config-client

	cp /data/temp-images/server.service /usr/lib/systemd/system/wireguard-config-server.service
	cp /data/temp-images/client@.service /usr/lib/systemd/system/wireguard-config-client@.service

	systemctl daemon-reload
	systemctl enable wireguard-config-server wireguard-config-client@service wireguard-config-client@normal

	nohup systemctl restart wireguard-config-server wireguard-config-client@service wireguard-config-client@normal &>/dev/null &
EOF
