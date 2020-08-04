//go:generate go-generate-struct-interface

package sharedConfig

type ConnectionOptions struct {
	Server        string `short:"s" long:"server" description:"config server ip:port" required:"true" env:"WIREGUARD_SERVER"`
	Password      string `short:"P" long:"password" description:"password for rpc calls" env:"WIREGUARD_PASSWORD"`
	GrpcInsecure  bool   `long:"insecure" description:"do not check server key (extreamly dangerous)" env:"WIREGUARD_TLS_INSECURE"`
	GrpcHostname  string `long:"server-name" description:"server hostname to verify with TLS" env:"WIREGUARD_TLS_SERVERNAME"`
	GrpcServerKey string `long:"server-ca" description:"use self-signed CA cert file" env:"WIREGUARD_TLS_CACERT"`
}
