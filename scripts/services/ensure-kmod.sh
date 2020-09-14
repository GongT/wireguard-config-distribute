#!/usr/bin/env bash

set -Eeuo pipefail

MOD="$1"

function wait_load() {
	while ! lsmod 2>/dev/null | grep -q --fixed-strings "$MOD" ; do
		echo "wait for kernel module $MOD to loaded..."
		sleep 1
	done
	echo "kernel module $MOD loaded."
}

function do_load() {
	echo "load kernel module $MOD"
	modprobe "$MOD" 2>/dev/null
	wait_load
}

if command -v systemd-detect-virt &>/dev/null ; then
	if systemd-detect-virt --container --quiet ; then
		wait_load
	else
		do_load
	fi
else
	echo "systemd-detect-virt not found..."
	do_load || true
fi
