//go:generate go-generate-struct-interface

package main

type serverProgramOptions struct {
	ServerName string `short:"n" long:"server-name" description:"(friendly) name of this server" default-mask:"$HOSTNAME" env:"HOSTNAME"`
	Password   string `short:"P" long:"password" description:"password for configure tool" default-mask:"generate one and save to storage" env:"WIREGUARD_PASSWORD"`

	ListenPath string `short:"u" long:"unix" description:"listen unix socket (and disable TCP)" env:"WIREGUARD_UNIX"`

	ListenPort      uint16 `short:"p" long:"port" description:"wireguard listening port" default:"51820" env:"WIREGUARD_PORT"`
	AutoForwardUpnp bool   `long:"upnp" description:"automantic open TCP port with UPnP/NAT-PMP" env:"WIREGUARD_UPNP"`

	PublicIp []string `long:"ip" description:"manually set public ip address of this device, disable auto detect" env:"WIREGUARD_PUBLIC_IP"`
	// IpUpnpDsiable bool   `long:"ip-noupnp" description:"disable detect public ipv4 by UPnP/NAT-PMP" env:"WIREGUARD_PUBLIC_IP_NO_UPNP"`
	IpHttpDsiable bool `long:"ip-nohttp" description:"disable detect public ip by request a http api" env:"WIREGUARD_PUBLIC_IP_NO_HTTP"`

	StorageLocation string `shourt:"s" long:"storage" description:"where to save data" default-mask:"~/.wireguard-config-server" env:"WIREGUARD_STORAGE"`

	GrpcInsecure  bool   `long:"insecure" description:"auto create self signed TLS key" env:"WIREGUARD_TLS_INSECURE"`
	GrpcServerKey string `long:"tls-keyfile" description:"use this TLS private key" env:"WIREGUARD_TLS_KEYFILE"`
	GrpcServerPub string `long:"tls-pubfile" description:"use this TLS public key" env:"WIREGUARD_TLS_PUBFILE"`

	DebugMode bool `long:"debug" short:"D" description:"enable debug mode" env:"WIREGUARD_CONFIG_DEVELOPMENT"`
}
