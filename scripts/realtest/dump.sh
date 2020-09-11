#!/usr/bin/env bash

set -Eeuo pipefail

set -a
source /etc/wireguard/client.conf
set +a

./scripts/run.ps1 tool dump 2>&1 \
	| sed 's/grpc connect ok\./\x1Bc/g; s/^==.*==$/\x1B[38;5;13;1m\0\x1B[0m/g; s/^> /🟢  /g'
