#!/usr/bin/env bash

set -Eeuo pipefail

cat <<- 'EOF' > /etc/wireguard/client.conf
	WIREGUARD_SERVER="grpc.services.gongt.me:443"
	WIREGUARD_TITLE="工作机(linux)"
	WIREGUARD_IPV6="true"
	WIREGUARD_PUBLIC_IP_NO_UPNP="true"
	WIREGUARD_PUBLIC_IP_NO_HTTP="true"
	WIREGUARD_NO_UPNP="true"
	WIREGUARD_CONFIG_DEVELOPMENT="true"
	WIREGUARD_REQUEST_IP="1.0"
EOF

set -a
source /etc/wireguard/client.conf
set +a

./scripts/run.ps1 client
