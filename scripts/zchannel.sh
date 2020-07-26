#!/usr/bin/env bash

set -Eeuo pipefail

function find_bridge_ip() {
	podman network inspect podman | grep -oE '"gateway": ".+",?$' | sed 's/"gateway": "\(.*\)".*/\1/g'
}

readonly ENVOY_PORT=$1
readonly GRPC_ADDR=$(find_bridge_ip)
readonly GRPC_PORT=51820
readonly ASSETS_PORT=9090
