#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..

export GOOS="linux"
export GOARCH="amd64"
export RHOST="router.home.gongt.me"

pwsh scripts/build.ps1 client

echo "Build with alpine..."
podman run --rm --workdir /data --volume "$(pwd):/data" --volume "$GOPATH:/go" "golang:alpine" go build -o ./dist/client.alpine ./cmd/wireguard-config-client

echo
echo

rsync dist/client.alpine scripts/services/client.init.sh $RHOST:/data/temp

cat <<- 'EOF' | ssh $RHOST bash
	set -x
	/etc/init.d/wireguard-config-client stop

	rm -f /usr/libexec/wireguard-config-client
	cp /data/temp/client.alpine /usr/libexec/wireguard-config-client
	cp /data/temp/client.init.sh /etc/init.d/wireguard-config-client

	/etc/init.d/wireguard-config-client enable
	/etc/init.d/wireguard-config-client start
EOF
