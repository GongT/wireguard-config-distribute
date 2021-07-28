#!/usr/bin/env bash

set -Eeuo pipefail

function die() {
	echo "$*" >&2
	exit 1
}

echo "stop systemd server"
systemctl stop wireguard-config-server.service

echo "running..."
export STATE_DIRECTORY=/var/lib/wireguard-config-server
./dist/server
