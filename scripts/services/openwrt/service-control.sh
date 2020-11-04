#!/usr/bin/env bash

set -Eeuo pipefail

declare -r ACTION=$1
/etc/init.d/wireguard-config-client "$ACTION"
