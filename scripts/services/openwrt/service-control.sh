#!/usr/bin/env bash

set -Eeuo pipefail

declare -r ACTION=$1
echo "[service-control] $ACTION ..." >&2
"/etc/init.d/wireguard-config" "$ACTION"
