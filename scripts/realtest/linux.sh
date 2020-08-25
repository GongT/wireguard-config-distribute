#!/usr/bin/env bash

set -Eeuo pipefail

set -a
source /etc/wireguard/client.normal.conf
source /etc/wireguard/client.conf
set +a

./scripts/run.ps1 client --group normal
