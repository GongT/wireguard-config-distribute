#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"

pwsh scripts/build.ps1 client

set -x
cp scripts/services/client@.service /usr/lib/systemd/system/wireguard-config-client@.service
systemctl daemon-reload
systemctl enable wireguard-config-client@normal
systemctl stop wireguard-config-client@normal
cp dist/client /usr/local/bin/wireguard-config-client
systemctl restart wireguard-config-client@normal