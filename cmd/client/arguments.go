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

	/* Public IPv4 */
	Ipv4Only bool `short:"4" long:"ipv4only" description:"disable connect anything with ipv6" env:"WIREGUARD_IPV4"`
	Ipv6Only bool `short:"6" long:"ipv6only" description:"disable ipv4 external ip auto detect" env:"WIREGUARD_IPV6"`

	PublicIp      string `long:"external-ip" description:"manually set public ipv4 address of this device, disable auto detect" env:"WIREGUARD_PUBLIC_IP"`
	PublicIp6     string `long:"external-ip6" description:"manually set public ipv6 address of this device, disable auto detect" env:"WIREGUARD_PUBLIC_IP6"`
	PublicPort    uint16 `long:"external-port" description:"manually set public port, if you are behind NAT device" default-mask:"UPnP or same with --port" env:"WIREGUARD_PUBLIC_PORT"`
	IpUpnpDsiable bool   `long:"external-ip-noupnp" description:"disable detect public ipv4 by UPnP/NAT-PMP" env:"WIREGUARD_PUBLIC_IP_NO_UPNP"`
	IpHttpDsiable bool   `long:"external-ip-nohttp" description:"disable detect public ip by request a http api" env:"WIREGUARD_PUBLIC_IP_NO_HTTP"`

	NoAutoForwardUpnp bool `long:"no-upnp-forward" description:"don't open port with UPnP/NAT-PMP" env:"WIREGUARD_NO_UPNP"`

	/* Local IPv4 */
	InternalIp string `long:"internal-ip" description:"manually set internal ipv4 address of this device" default-mask:"auto detect" env:"WIREGUARD_PRIVATE_IP"`

	/* Debug */
	DryRun bool `long:"dry" description:"do not create any interface" env:"WIREGUARD_DRY_RUN"`
}
