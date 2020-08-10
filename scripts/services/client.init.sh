#!/bin/bash /etc/rc.common

source /lib/functions/network.sh

USE_PROCD=1
START=80

BIN=/usr/libexec/wireguard-config-client

create_pid() {
	echo "/tmp/run/wireguard-config-client-$_SECTION.pid"
}

declare -ar CONFIGS=(server password tls_insecure tls_servername tls_cacert port mtu network request_ip title hosts_file public_port dry_run log config_development)

_config_contains() {
	local name
	for name in "${CONFIGS[@]}"; do
		if [[ "$name" == "$1" ]]; then
			return 0
		fi
	done
	return 1
}

x() {
	echo "$*"
	"$@"
}

create_instance() {
	echo " * $_SECTION" >&2

	local debug_file="/tmp/start.wireguard.$_SECTION.sh"
	echo "" > "$debug_file"

	procd_open_instance "$_SECTION"

	if [[ ! "${CONFIGS[port]}" ]]; then
		echo "Config section $_SECTION must have a port option."
		exit 1
	fi

	### environment
	local name EnvList=("WIREGUARD_GROUP=$_SECTION")
	for name in "${!_CONFIGS[@]}"; do
		echo "export WIREGUARD_${name^^}='${_CONFIGS[$name]}'" >> "$debug_file"
		EnvList+=("WIREGUARD_${name^^}=${_CONFIGS[$name]}")
	done
	procd_set_param env "${EnvList[@]}"
	echo "exec ${RUN[*]}" >> "$debug_file"

	### cmd
	local RUN=(
		"$BIN"
		--external-ip-noupnp
		--external-ip-nohttp
		--no-upnp-forward
		"--hostname=$HOSTNAME"
		"--external-ip=$WAN_IP"
		"--external-ip6=$WAN_IP6"
		"--internal-ip=$LAN_IP"
	)
	procd_set_param command "${RUN[@]}"

	### configs
	procd_set_param respawn "${respawn_threshold:-3600}" "${respawn_timeout:-5}" "${respawn_retry:-2}"
	procd_set_param file /etc/config/wireguard_config
	procd_set_param stdout 1
	procd_set_param stderr 1
	procd_set_param user root
	procd_set_param pidfile "$(create_pid)"
	procd_close_instance
}

config_cb() {
	if [[ -n "$_SECTION" ]]; then
		# echo "CREATE!!! $WIREGUARD_INTERFACE_NAME"
		create_instance "$_SECTION"
		_SECTION=
		_CONFIGS=()
	fi
	if [[ "$1" == "interface" ]]; then
		_SECTION="$2"
	elif [[ "$1" ]]; then
		echo "Warn: Unknown section: $1 $2" >&2
	fi
}

option_cb() {
	if [[ -z "$_SECTION" ]]; then
		return
	fi
	local KEY=$1 VALUE=$2
	# echo "option_cb set $1=$2"
	if _config_contains "$KEY"; then
		_CONFIGS[$KEY]=$VALUE
	else
		echo "Warn: invalid option: $KEY" >&2
	fi
}

start_service() {
	local name _SECTION "${CONFIGS[@]}"
	local -A _CONFIGS=()
	echo "starting..." >&2

	local LAN_IP

	LAN_IP=$(uci get network.lan.ipaddr)
	if [[ ! "$LAN_IP" ]]; then
		echo "Failed find LAN IP" >&2
		exit 1
	fi

	local WAN_IP WAN_IP6 NET_IF
	network_find_wan NET_IF
	network_get_ipaddr WAN_IP "${NET_IF}"
	if [[ ! "$WAN_IP" ]]; then
		echo "Failed find WAN IPv4" >&2
		exit 1
	fi

	network_find_wan6 NET_IF
	if [[ "$NET_IF" ]]; then
		network_get_ipaddr6 WAN_IP6 "${NET_IF}"
		if [[ ! "$WAN_IP6" ]]; then
			echo "Failed find WAN IPv6 address" >&2
		fi
	else
		echo "Failed find WAN IPv6 interface" >&2
	fi

	config_load wireguard_config
}
