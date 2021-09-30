#!/bin/bash /etc/rc.common

source /lib/functions/network.sh

USE_PROCD=1
START=80

die() {
	echo "$*" >&2
	exit 1
}

create_pid() {
	echo "/tmp/run/wireguard-config-$_SECTION.pid"
}
create_bin() {
	echo "/usr/local/libexec/wireguard-config/${1:-"${_SECTION_KIND}.alpine"}"
}

update_now() {
	die "No binary ($(create_bin))"
}

reset_config() {
	if [[ $1 == client ]]; then
		_CONFIGS=(
			[EXTIP_NATIVE]=true
			[EXTIP4_NO_UPNP]=true
			[EXTIP4_API]=false
			[EXTIP6_API]=false
			[NO_UPNP]=true
			[PRIVATE_IP]="$LAN_IP"
		)
	else
		_CONFIGS=()
	fi
}
declare -a KNOWN_CONFIGS=()
check_binary() {
	local BIN
	BIN=$(create_bin)
	"$BIN" --version &>/dev/null || update_now
	"$BIN" --version &>/dev/null || die "Binary file not exists or not executable, please check environment."

	mapfile -t KNOWN_CONFIGS < <("$BIN" --help | grep -Eo '\[\$WIREGUARD_.*]' | sed 's/^\[\$WIREGUARD_//' | sed 's/]$//')
}

_config_contains() {
	local name
	for name in "${KNOWN_CONFIGS[@]}"; do
		if [[ $name == "${1^^}" ]]; then
			return 0
		fi
	done
	return 1
}

x() {
	echo "$*"
	"$@"
}

finalize_instance() {
	{
		echo " * $_SECTION_KIND - $_SECTION ($debug_file)"
		echo "   $(create_bin)"
		for name in "${!_CONFIGS[@]}"; do
			echo "   - WIREGUARD_${name^^}=${_CONFIGS[$name]}"
		done
	} >&2

	procd_open_instance "$_SECTION"

	# if [[ ! ${_CONFIGS[PORT]} ]]; then
	# 	echo "Config section $_SECTION must have a port option."
	# 	exit 1
	# fi

	### environment
	local name EnvList=("HOSTNAME=$(uci get system.@system[0].hostname)" "WIREGUARD_GROUP=$_SECTION")
	for name in "${!_CONFIGS[@]}"; do
		EnvList+=("WIREGUARD_${name^^}=${_CONFIGS[$name]}")
	done
	procd_set_param env "${EnvList[@]}"

	### cmd
	procd_set_param command "$(create_bin)"

	### configs
	procd_set_param respawn "${respawn_threshold:-3600}" "${respawn_timeout:-5}" "${respawn_retry:-2}"
	procd_set_param file /etc/config/wireguard_config # TODO
	procd_set_param stdout 1
	procd_set_param stderr 1
	procd_set_param user root
	procd_set_param pidfile "$(create_pid)"
	procd_close_instance
}

config_cb() {
	if [[ -n $_SECTION ]]; then
		# echo "CREATE!!! $WIREGUARD_INTERFACE_NAME"
		finalize_instance "$_SECTION"
		_SECTION=
		_SECTION_KIND=
		reset_config "$1"
	fi
	if [[ $1 == "client" ]] || [[ $1 == "server" ]]; then
		_SECTION_KIND="$1"
		_SECTION="$2"
		check_binary
	elif [[ "$1" ]]; then
		echo "Warn: Unknown section: $1 $2" >&2
		return
	fi
}

option_cb() {
	if [[ -z $_SECTION ]]; then
		return
	fi
	local KEY=$1 VALUE=$2
	# echo "option_cb set $1=$2"
	if _config_contains "$KEY"; then
		_CONFIGS[${KEY^^}]=$VALUE
	else
		echo "Warn: invalid option: $KEY" >&2
	fi
}

start_service() {
	local name _SECTION
	local -A _CONFIGS
	echo "starting..." >&2

	local LAN_IP
	LAN_IP=$(uci get network.lan.ipaddr)
	if [[ ! $LAN_IP ]]; then
		echo "Failed find LAN IP" >&2
		exit 1
	fi

	reset_config
	config_load wireguard_config
}
