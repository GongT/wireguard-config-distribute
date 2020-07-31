#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

echo -ne "\ec"

export GOOS="linux"
export GOARCH="amd64"
export RHOST="services.gongt.me"

pwsh scripts/build.ps1 server

rsync -v dist/server $RHOST:/tmp/wireguard-config-server

ssh $RHOST /tmp/wireguard-config-server -D --insecure --unix /dev/shm/container-shared-socksets/grpc.wireguard.sock
