//go:generate go-generate-struct-interface

package config

type internalOptions struct {
	StandardOutputPath string `long:"pipe-output"`
	IsElevated         bool   `long:"is-elevate"`
}
type commonOptions struct {
	DebugMode   bool   `long:"debug" short:"D" description:"enable debug mode" env:"WIREGUARD_CONFIG_DEVELOPMENT"`
	LogFilePath string `long:"logfile" description:"save output to file" env:"WIREGUARD_LOG"`
	ShowVersion bool   `long:"version" short:"V" description:"show program version and exit"`
}
