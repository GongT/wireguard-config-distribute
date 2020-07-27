//go:generate go run ../../tools/generate.go

package main

type downloadCaAction struct {
	Output string `short:"O" long:"output" description:"output file path" default:"./server.cert"`
}

type toolProgramOptions struct {
	Server   string `short:"s" long:"server" description:"config server ip:port" default:"127.0.0.1:51820" env:"WIREGUARD_SERVER"`
	Password string `short:"P" long:"password" description:"server password (required when connect to remote server)" env:"WIREGUARD_PASSWORD"`

	DownloadCA downloadCaAction `command:"download-ca" alias:"auth" description:"Download server's self-signed CA cert file"`

	DebugMode bool `long:"debug" short:"D" description:"enable debug mode" env:"WIREGUARD_CONFIG_DEVELOPMENT"`
}
