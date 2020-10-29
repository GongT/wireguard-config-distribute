#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_FILE=$1
shift

declare -a EXTRA_FILES=()

declare -r TMPF="$(mktemp)"
trap 'cd / ; rm -f "$TMPF"' EXIT

tar -czf "$TMPF" \
	--absolute-names \
	"--transform=s#$SCRIPT_FILE#__script.sh#" \
	"${EXTRA_FILES[@]}" \
	"$@" \
	"$SCRIPT_FILE"

du --bytes "$TMPF" | cut -f1

cat "$TMPF"

if ! [[ -t 0 ]]; then
	cat
fi
