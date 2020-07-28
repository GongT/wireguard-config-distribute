//go:generate go run ../../tools/generate.go

package main

type clientProgramOptions struct {
	/* wg interface */
	ListenPort uint16 `short:"p" long:"port" description:"wireguard listening port" default:"51820" env:"WIREGUARD_PORT"`
	MTU        uint16 `long:"mtu" description:"wireguard interface MTU" env:"WIREGUARD_MTU"`

	/* config server and self config */
	Server      string `short:"s" long:"server" description:"config server ip:port" required:"true" env:"WIREGUARD_SERVER"`
	NetworkName string `short:"n" long:"netgroup" description:"a (friendly) name of local network, all machines in one local network should same" default-mask:"no default network" env:"WIREGUARD_NETWORK"`
	JoinGroup   string `short:"g" long:"group" description:"join which VPN network" default:"default" env:"WIREGUARD_GROUP"`
	PerferIp    string `long:"perfer-ip" description:"request to use this VPN ip, only last two digist" default-mask:"allocate by server" env:"WIREGUARD_REQUEST_IP"`

	/* Self application config */
	Title    string `short:"t" long:"title" description:"human readable name of this machine" default-mask:"same with hostname" env:"WIREGUARD_TITLE"`
	Hostname string `long:"hostname" description:"custom hostname (insteadof environment)" default-mask:"use HOSTNAME" env:"HOSTNAME"`
	HostFile string `long:"hosts-file" description:"watch and read hosts file" default:"/etc/hosts" env:"WIREGUARD_HOSTS_FILE"`

	/* Public IPv4 */
	// Ipv4Only bool `short:"4" long:"ipv4only" description:"disable incoming ipv6 connection" env:"WIREGUARD_IPV4"`
	Ipv6Only bool `short:"6" long:"ipv6only" description:"disable incoming ipv4 connection, disable ipv4 external ip auto detect" env:"WIREGUARD_IPV6"`

	PublicIp      string `long:"external-ip" description:"manually set public ipv4 address of this device, disable auto detect" env:"WIREGUARD_PUBLIC_IP"`
	PublicIp6     string `long:"external-ip6" description:"manually set public ipv6 address of this device, disable auto detect" env:"WIREGUARD_PUBLIC_IP6"`
	IpUpnpDsiable bool   `long:"external-ip-noupnp" description:"disable detect public ipv4 by UPnP/NAT-PMP" env:"WIREGUARD_PUBLIC_IP_NO_UPNP"`
	IpHttpDsiable bool   `long:"external-ip-nohttp" description:"disable detect public ip by request a http api" env:"WIREGUARD_PUBLIC_IP_NO_HTTP"`

	NoAutoForwardUpnp bool `long:"no-upnp-forward" description:"don't open port with UPnP/NAT-PMP" env:"WIREGUARD_NO_UPNP"`

	/* Local IPv4 */
	InternalIp []string `long:"internal-ip" description:"manually set internal ipv4 address of this device" default-mask:"auto detect" env:"WIREGUARD_PRIVATE_IP"`

	/* grpc endpoint */
	GrpcInsecure  bool   `long:"insecure" description:"do not check server key (extreamly dangerous)" env:"WIREGUARD_TLS_INSECURE"`
	GrpcHostname  string `long:"server-name" description:"server hostname to verify with TLS" env:"WIREGUARD_TLS_SERVERNAME"`
	GrpcServerKey string `long:"server-ca" description:"use self-signed CA cert file" env:"WIREGUARD_TLS_CACERT"`

	/* debug */
	DebugMode bool `long:"debug" short:"D" description:"enable debug mode" env:"WIREGUARD_CONFIG_DEVELOPMENT"`
}
