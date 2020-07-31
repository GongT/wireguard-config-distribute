#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"
export RHOST="services.gongt.me"

pwsh scripts/build.ps1 server

scp dist/server $RHOST:/usr/local/bin/wireguard-config-server
scp scripts/services/server.service $RHOST:/usr/lib/systemd/system/wireguard-config-server.service
