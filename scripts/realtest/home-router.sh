#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

echo -ne "\ec"

export GOOS="linux"
export GOARCH="amd64"
export RHOST="router.home.gongt.me"

pwsh scripts/build.ps1 client

rsync --progress --verbose --exclude dist . $RHOST:/data/temp/wireguard-config-client-test

cat <<- 'EOF' | tee | ssh $RHOST bash
	mkdir -p /etc/wireguard
	cat <<- 'EEOF' > /etc/wireguard/client.conf
		WIREGUARD_SERVER="grpc.gateway.gongt.me:443"
		WIREGUARD_NETWORK="home"
		WIREGUARD_TITLE="路由器"
		WIREGUARD_PUBLIC_IP_NO_UPNP="true"
		WIREGUARD_PUBLIC_IP_NO_HTTP="true"
		WIREGUARD_NO_UPNP="true"
		WIREGUARD_CONFIG_DEVELOPMENT="true"
		WIREGUARD_REQUEST_IP="0.0"
	EEOF

	set -a
	source /etc/wireguard/client.conf
	set +a

	cd /data/temp/wireguard-config-client-test
	go run ./cmd/wireguard-config-client
EOF
