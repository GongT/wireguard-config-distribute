#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"

pwsh scripts/build.ps1 client

scp dist/client router.home.gongt.me:/usr/local/bin/wireguard-config-client
scp scripts/services/client.init.sh router.home.gongt.me:/etc/init.d/wireguard-config-client

ssh router.home.gongt.me /etc/init.d/wireguard-config-client 
