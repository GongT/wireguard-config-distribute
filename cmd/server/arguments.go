//go:generate go run ../../tools/generate.go

package main

type serverProgramOptions struct {
	ServerName string `short:"n" long:"server-name" description:"(friendly) name of this server" env:"WIREGUARD_SERVER_NAME"`

	ListenPath string `short:"u" long:"unix" description:"listen unix socket (and disable TCP)" env:"WIREGUARD_UNIX"`

	ListenPort      uint16 `short:"p" long:"port" description:"wireguard listening port" default:"51820" env:"WIREGUARD_PORT"`
	AutoForwardUpnp bool   `long:"upnp" description:"automantic open TCP port with UPnP/NAT-PMP" env:"WIREGUARD_UPNP"`

	StorageLocation string `shourt:"s" long:"storage" description:"where to save data" default-mask:"~/.wireguard-config-server" env:"WIREGUARD_STORAGE"`

	GrpcInsecure  bool   `long:"insecure" description:"auto create self signed TLS key" env:"WIREGUARD_TLS_INSECURE"`
	GrpcServerKey string `long:"tls-keyfile" description:"use this TLS private key" env:"WIREGUARD_TLS_KEYFILE"`
	GrpcServerPub string `long:"tls-pubfile" description:"use this TLS public key" env:"WIREGUARD_TLS_PUBFILE"`

	DebugMode bool `long:"debug" short:"D" description:"enable debug mode" env:"WIREGUARD_CONFIG_DEVELOPMENT"`
}
