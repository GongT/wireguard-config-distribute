#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"
export RHOST="services.gongt.me"

pwsh scripts/build.ps1 server

rsync dist/server $RHOST:/usr/local/bin/wireguard-config-server.update
scp scripts/services/server.service $RHOST:/usr/lib/systemd/system/wireguard-config-server.service

cat <<- 'EOF' | ssh $RHOST bash 
	set -xEeuo pipefail
	systemctl daemon-reload
	systemctl enable wireguard-config-server
	systemctl stop wireguard-config-server
	rm -f /usr/local/bin/wireguard-config-server
	mv /usr/local/bin/wireguard-config-server.update /usr/local/bin/wireguard-config-server
	nohup systemctl restart wireguard-config-server &>/dev/null &
EOF
