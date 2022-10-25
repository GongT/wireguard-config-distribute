//go:generate go-generate-struct-interface

package main

import "github.com/gongt/wireguard-config-distribute/internal/client/sharedConfig"

type clientProgramOptionsBase struct {
	ConnectionOptions sharedConfig.ConnectionOptions `group:"Connection Options"`

	/* wg interface */
	ListenPort    uint16 `short:"p" long:"port" description:"wireguard listening port" default-mask:"random select" env:"WIREGUARD_PORT"`
	InterfaceName string `short:"i" long:"interface" description:"wireguard interface name (must not exists)" default-mask:"wg_${group}" env:"WIREGUARD_INTERFACE_NAME"`
	MTU           uint16 `long:"mtu" description:"wireguard interface MTU" env:"WIREGUARD_MTU"`

	/* config server and self config */
	NetworkName string `short:"n" long:"netgroup" description:"a (friendly) name of local network, all machines in one local network should same" env:"WIREGUARD_NETWORK"`
	JoinGroup   string `short:"g" long:"group" description:"join which VPN network" default:"default" env:"WIREGUARD_GROUP"`
	PerferIp    string `long:"perfer-ip" description:"request to use this VPN ip, only last two digist" default-mask:"allocate by server" env:"WIREGUARD_REQUEST_IP"`

	/* Self application config */
	Title     string `short:"t" long:"title" description:"human readable name of this machine" default-mask:"same with hostname" env:"WIREGUARD_TITLE"`
	Hostname  string `long:"hostname" description:"custom hostname (insteadof environment)"`
	HostFile  string `long:"hosts-file" description:"watch and read hosts file" default:"/etc/hosts" env:"WIREGUARD_HOSTS_FILE"`
	MachineID string `long:"machine-id" description:"global unique id of this machine" env:"WIREGUARD_MACHINE_ID"`

	/* Debug */
	DryRun bool `long:"dry" description:"do not create any interface" env:"WIREGUARD_DRY_RUN"`
}

type notMoveArguments struct {
	PublicPort uint16 `long:"external-port" description:"manually set public port, if you are behind NAT device" default-mask:"UPnP or same with --port" env:"WIREGUARD_PUBLIC_PORT"`

	Movable bool `long:"moveable" description:"disable connect from local network" env:"WIREGUARD_MOVABLE"`

	NoPublicNetwork bool     `long:"disable-listen" description:"mark there is no way to access this device from internet" env:"WIREGUARD_PRIVATE"`
	PublicIp        []string `long:"ip" description:"manually set public ip address of this device" env:"WIREGUARD_EXTIP"`
	Gateway         bool     `long:"ip-native" description:"read external ip from system network card (and disable all detect methods below)" env:"WIREGUARD_EXTIP_NATIVE"`
	IpUpnpDisable   bool     `long:"ip-no-upnp" description:"try UPnP to get ipv4" env:"WIREGUARD_EXTIP_NO_UPNP"`
	IpApi          string   `long:"ip-api" description:"request this url to get ipv4 (disable when set to empty)" default:"http://show-my-ip.gongt.me" env:"WIREGUARD_EXTIP_API"`

	NoAutoForwardUpnp bool `long:"no-upnp-forward" description:"don't open port with UPnP/NAT-PMP" env:"WIREGUARD_NO_UPNP"`

	/* Local IPv4 */
	InternalIp string `long:"internal-ip" description:"manually set local ipv4 address of this device" default-mask:"auto detect" env:"WIREGUARD_PRIVATE_IP"`
}
