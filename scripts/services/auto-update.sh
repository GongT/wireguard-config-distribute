#!/usr/bin/env bash

set -Eeuo pipefail
cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

function die() {
	echo -e "\e[38;5;9m$*\e[0m"
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

function getCurrentInstalledVersion() {
	local OUTPUT
	OUTPUT=$("$BINARY_FILE" --version 2>/dev/null || true)
	if [[ ! $OUTPUT ]]; then
		return
	fi

	echo "$OUTPUT" | grep -o 'Hash:.*' | awk '{print $2}'
}

if ncat -z --wait 1 proxy-server. 3271 &>/dev/null; then
	export PROXY="http://proxy-server.:3271/"
	export https_proxy=${PROXY} http_proxy=${PROXY} all_proxy=${PROXY} HTTPS_PROXY=${PROXY} HTTP_PROXY=${PROXY} ALL_PROXY=${PROXY} NO_PROXY="10.*,192.*,127.*,172.*"
	info "using proxy server $PROXY"
fi

export PROJ_TYPE="${1}"

BIN_TYPE=''
if ldd /proc/$$/exe | grep -q -- 'musl'; then
	BIN_TYPE='.alpine'
fi

declare -r REPO='GongT/wireguard-config-distribute'
declare -r GET_FILE="$PROJ_TYPE$BIN_TYPE"
declare -r DIST_ROOT="/usr/local/libexec/wireguard-config"
declare -r BINARY_FILE="$DIST_ROOT/$GET_FILE"
declare -r LATEST_URL="https://api.github.com/repos/$REPO/releases?page=1&per_page=1"

mkdir -p "$DIST_ROOT"

info "检查 $REPO 版本……"

if [[ -e $BINARY_FILE ]]; then
	declare -r VERSION_LOCAL=$(getCurrentInstalledVersion)
else
	declare -r VERSION_LOCAL=
fi
mute "    本地版本： $VERSION_LOCAL"

GITHUB_AUTH=()
if [[ -e ~/.github-token ]]; then
	GITHUB_AUTH=(--header "authorization: Bearer $(<~/.github-token)")
fi


mute "    来源： $LATEST_URL"
declare RELEASE_DATA
RELEASE_DATA=$(curl -s "${GITHUB_AUTH[@]}" "$LATEST_URL")
RELEASE_DATA=$(echo "$RELEASE_DATA" | jq -M -c ".[0] // null") || {
	curl -v "$LATEST_URL" || true
	die "无法获取最新版本信息"
}
if [[ $RELEASE_DATA == "null" ]]; then
	die "failed get release data."
fi

declare REMOTE_VERSION=$(echo "$RELEASE_DATA" | jq -r -M -c ".target_commitish")
mute "    远程版本： $REMOTE_VERSION"

if [[ $VERSION_LOCAL == $REMOTE_VERSION ]]; then
	info " * 已是最新版本"
	mute "    文件:   $BINARY_FILE"
	[[ ${DISABLE_RESTART:-} ]] || bash "service-control.sh" start
	exit 0
fi

downloadUrl=$(buildGithubReleaseUrl "$(echo "$RELEASE_DATA" | jq -r -M -c ".tag_name")")
info " * 有更新，开始下载："
mute "    远程: $downloadUrl"
mute "    本地:   $BINARY_FILE"
wget --quiet --show-progress --progress=bar -O "$BINARY_FILE.downloading" "$downloadUrl"
echo >&2
rm -f "$BINARY_FILE"
mv "$BINARY_FILE.downloading" "$BINARY_FILE"
chmod a+x "$BINARY_FILE"

echo -n "当前版本："
getCurrentInstalledVersion

echo
echo "Ah, that's ♂ good."

[[ ${DISABLE_RESTART:-} ]] || bash "service-control.sh" restart
