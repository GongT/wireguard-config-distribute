#!/usr/bin/env bash

set -Eeuo pipefail

declare -r PACKAGE_SIZE EXTRA_INPUT
declare -r TMPD=$(mktemp -d)
trap 'cd / ; rm -rf "$TMPD"' EXIT

cd "$TMPD"

if [[ $PACKAGE_SIZE -le 0 ]]; then
	echo "Invalid header" >&2
	exit 1
fi

function run() {
	echo "receive $PACKAGE_SIZE bytes package"
	dd if=/proc/$$/fd/0 count=1 bs=$PACKAGE_SIZE status=none \
		| tar -xzf- --transform "s#^/#./#"

	chmod a+x __script.sh

	# echo "execute remote script with args: [$#] $*"
	bash ./__script.sh "$@"

	echo "remote script done." >&2

	exit 0
}

run "$@"
