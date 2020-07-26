//go:generate go run ../../tools/generate.go

package main

type downloadCaAction struct {
	Output string `short:"O" long:"output" description:"output file path" default:"./server.cert"`
}

type toolProgramOptions struct {
	Server   string `short:"s" long:"server" description:"config server ip:port" required:"true" env:"WIREGUARD_SERVER"`
	Password string `short:"P" long:"password" description:"server password" required:"true" env:"WIREGUARD_PASSWORD"`

	DownloadCA downloadCaAction `command:"download-ca" alias:"auth" description:"Download server's self-signed CA cert file"`

	DebugMode bool `long:"debug" short:"D" description:"enable debug mode" env:"WIREGUARD_CONFIG_DEVELOPMENT"`
}
