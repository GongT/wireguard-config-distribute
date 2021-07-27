#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_FILE=$1
shift

MY_PATH="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
declare -a EXTRA_FILES=()

declare -r TMPF="$(mktemp)"
trap 'cd / ; rm -f "$TMPF"' EXIT

tar -czf "$TMPF" \
	--absolute-names \
	"--transform=s#$SCRIPT_FILE#__script.sh#" \
	"${EXTRA_FILES[@]}" \
	"$@" \
	"$SCRIPT_FILE"

PACKAGE_SIZE=$(du --bytes "$TMPF" | cut -f1)
declare -p PACKAGE_SIZE

EXTRA_INPUT=no
if ! [[ -t 0 ]]; then
	EXTRA_INPUT=yes
fi
declare -p EXTRA_INPUT

echo "send $PACKAGE_SIZE bytes package" >&2

cat "$MY_PATH/script-receiver.sh"
# echo "extracter sent" >&2
cat "$TMPF"
# echo "package sent" >&2

if ! [[ -t 0 ]]; then
	# echo "redirect input..." >&2
	cat
	# echo "redirect input done." >&2
fi
