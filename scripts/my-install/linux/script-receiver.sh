#!/usr/bin/env bash

set -Eeuo pipefail

declare -r TMPD=$(mktemp -d)
trap 'cd / ; rm -rf "$TMPD"' EXIT

cd "$TMPD"

SIZE=""
while true; do
	I=$(dd if=/dev/stdin count=1 bs=1 status=none)
	if [[ -z $I ]] || [[ $I == $'\n' ]]; then
		break
	fi
	SIZE+="$I"
done

if [[ $SIZE -le 0 ]]; then
	echo "Invalid header" >&2
	exit 1
fi

dd if=/dev/stdin count=1 bs=$SIZE status=none \
	| tar -xzf- \
		--transform "s#^/#./#"

chmod a+x __script.sh

./__script.sh "$@"
