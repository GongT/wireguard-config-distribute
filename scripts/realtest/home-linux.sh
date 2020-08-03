#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

echo -ne "\ec"

export GOOS="linux"
export GOARCH="amd64"
export RHOST="home.gongt.me"

pwsh scripts/build.ps1 client

rsync -v dist/client $RHOST:/tmp/wireguard-config-client

ssh -tt $RHOST bash -c "
	export WIREGUARD_SERVER='grpc.services.gongt.me:443'
	export WIREGUARD_TITLE='家里服务器主机'
	export WIREGUARD_CONFIG_DEVELOPMENT='true'
	export WIREGUARD_REQUEST_IP='0.10'
	exec /tmp/wireguard-config-client
"
