#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")/.."

declare -r TYPE=tool

# echo "Generate protobuf..."
# ./scripts/create-protobuf.ps1

echo "Build $TYPE..."
go build -o dist/$TYPE ./cmd/$TYPE

unset USERNAME

echo "Copy binary..."
scp dist/$TYPE home.gongt.me:/tmp

echo "Execute..."
ssh -tt home.gongt.me /tmp/$TYPE
