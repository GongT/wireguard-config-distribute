#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"
export RHOST="services.gongt.me"

pwsh scripts/build.ps1 server
pwsh scripts/build.ps1 client

rsync dist/server $RHOST:/usr/local/bin/wireguard-config-server.update
rsync dist/client $RHOST:/usr/local/bin/wireguard-config-client.update
scp scripts/services/server.service $RHOST:/usr/lib/systemd/system/wireguard-config-server.service
scp scripts/services/client.service $RHOST:/usr/lib/systemd/system/wireguard-config-client.service

cat <<- 'EOF' | ssh $RHOST bash 
	set -xEeuo pipefail
	systemctl daemon-reload
	systemctl enable wireguard-config-server wireguard-config-client
	systemctl stop wireguard-config-server wireguard-config-client
	rm -f /usr/local/bin/wireguard-config-server /usr/local/bin/wireguard-config-client
	mv /usr/local/bin/wireguard-config-server.update /usr/local/bin/wireguard-config-server
	mv /usr/local/bin/wireguard-config-client.update /usr/local/bin/wireguard-config-client
	nohup systemctl restart wireguard-config-server wireguard-config-client &>/dev/null &
EOF
