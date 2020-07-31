#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"
export RHOST="router.home.gongt.me"

pwsh scripts/build.ps1 client

rsync dist/client $RHOST:/usr/local/bin/wireguard-config-client
scp scripts/services/client.init.sh $RHOST:/etc/init.d/wireguard-config-client

ssh $RHOST /etc/init.d/wireguard-config-client enable
