#!/usr/bin/env bash

set -Eeuo pipefail
cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")/.."

export HOSTNAME=test
export WIREGUARD_CONFIG_DEVELOPMENT=true
export WIREGUARD_PUBLIC_IP_NO_HTTP=true
export WIREGUARD_STORAGE=/tmp/xxxaaaa
export WIREGUARD_PASSWORD=123456

go test -trace=trace.out -cpuprofile cpu.out ./cmd/wireguard-config-server/*.go -- xxx
