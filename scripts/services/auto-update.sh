#!/usr/bin/env bash

set -Eeuo pipefail
cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

function die() {
	echo -e "\e[38;5;9m$*\3[0m"
	exit 1
}
function info() {
	echo -e "$*" >&2
}
function mute() {
	echo -e "\e[2m$*\e[0m" >&2
}

function buildGithubReleaseUrl() {
	echo "https://github.com/$REPO/releases/download/$1/$GET_FILE"
}

function allServices() {
	local ACTION=$1

	set +e
	mapfile -t SERVICES < <(systemctl list-units --no-pager --no-legend wireguard-config-* | cut -d ' ' -f 1)
	info "$ACTION services:"
	for I in "${SERVICES[@]}"; do
		if [[ $I == "wireguard-config-auto-update"* ]]; then
			mute "  * $I ... (skip)"
			continue
		fi
		info "  * $I ..."
		systemctl --no-block "$ACTION" "$I"
	done
}

if getent hosts proxy-server. &>/dev/null && curl --proxy proxy-server.:3271 www.google.com &>/dev/null; then
	export PROXY="http://proxy-server.:3271/"
	export https_proxy=${PROXY} http_proxy=${PROXY} all_proxy=${PROXY} HTTPS_PROXY=${PROXY} HTTP_PROXY=${PROXY} ALL_PROXY=${PROXY} NO_PROXY="10.*,192.*,127.*,172.*"
	info "using proxy server $PROXY"
fi

PROJ_TYPE="${1}"

BIN_TYPE=''
if ldd /proc/$$/exe | grep -q -- 'musl'; then
	BIN_TYPE='.alpine'
fi

declare -r REPO='GongT/wireguard-config-distribute'
declare -r GET_FILE="$PROJ_TYPE$BIN_TYPE"
declare -r DIST_ROOT="/usr/local/libexec/wireguard-config-$PROJ_TYPE"
declare -r BINARY_FILE="$DIST_ROOT/$GET_FILE"
declare -r LATEST_URL="https://api.github.com/repos/$REPO/releases?page=1&per_page=1"

mkdir -p "$DIST_ROOT"

declare -r VERSION_FILE="$DIST_ROOT/$GET_FILE.version.txt"

info "检查 $REPO 版本……"
mute "    来源： $LATEST_URL"
declare RELEASE_DATA
RELEASE_DATA=$(curl -s "$LATEST_URL" | jq -M -c ".[0] // null")
if [[ $RELEASE_DATA == "null" ]]; then
	die "failed get release data."
fi

if [[ -e $VERSION_FILE ]]; then
	declare -ir VERSION_LOCAL=$(<"$VERSION_FILE")
else
	declare -ir VERSION_LOCAL=0
fi

declare -i REMOTE_VERSION=$(echo "$RELEASE_DATA" | jq -r -M -c ".id")
if [[ $VERSION_LOCAL -eq $REMOTE_VERSION ]]; then
	info " * 已是最新版本"
	mute "    文件:   $BINARY_FILE"
	allServices start
	exit 0
fi

downloadUrl=$(buildGithubReleaseUrl "$(echo "$RELEASE_DATA" | jq -r -M -c ".tag_name")")
info " * 有更新，开始下载："
mute "    远程: $downloadUrl"
mute "    本地:   $BINARY_FILE"
wget --quiet --show-progress --progress=bar -O "$BINARY_FILE.downloading" "$downloadUrl"
rm -f "$BINARY_FILE"
mv "$BINARY_FILE.downloading" "$BINARY_FILE"
chmod a+x "$BINARY_FILE"

echo "$REMOTE_VERSION" >"$VERSION_FILE"

echo -n "当前版本："
"$BINARY_FILE" --version 2>/dev/null || die "binary file not executable"

echo
echo "Ah, that's ♂ good."

allServices restart
